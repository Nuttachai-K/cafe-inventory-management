package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/jackc/pgx/v5"
)

type InventoryRepository interface {
	GetByID(ctx context.Context, id int) (*model.InventoryWithProduct, error)
	GetAll(ctx context.Context, productName string, limit int) ([]model.InventoryWithProduct, error)
	Create(ctx context.Context, inventory *model.Inventory) error
	UpdateStock(ctx context.Context, id int, quantity int, isAdjust bool) (int, int, error)
	Delete(ctx context.Context, id int) error
}

type inventoryRepository struct {
	db DBTX
}

var ErrInsufficientStock = errors.New("insufficient stock")

func NewInventoryRepository(db DBTX) InventoryRepository {
	return &inventoryRepository{
		db: db,
	}
}

func (i *inventoryRepository) GetByID(ctx context.Context, id int) (*model.InventoryWithProduct, error) {

	var inventory model.InventoryWithProduct
	err := i.db.QueryRow(
		ctx,
		`SELECT
			inventory.id,
			inventory.product_id AS product_id,
			products.name AS product_name,
			cafes.id AS cafe_id,
			cafes.name AS cafe_name,
			inventory.stock_quantity,
			inventory.updated_at
		FROM inventory
		JOIN products
		ON inventory.product_id = products.id
		JOIN cafes 
		ON products.cafe_id = cafes.id
		WHERE inventory.product_id = $1`,
		id,
	).Scan(
		&inventory.ID,
		&inventory.ProductId,
		&inventory.ProductName,
		&inventory.CafeId,
		&inventory.CafeName,
		&inventory.StockQuantity,
		&inventory.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &inventory, err
}

func (i *inventoryRepository) GetAll(ctx context.Context, productName string, limit int) ([]model.InventoryWithProduct, error) {

	query := `
	SELECT
		inventory.id,
		inventory.product_id AS product_id,
		products.name AS product_name,
		cafes.id AS cafe_id,
		cafes.name AS cafe_name,
		inventory.stock_quantity,
		inventory.updated_at
	FROM inventory
	JOIN products
	ON inventory.product_id = products.id
	JOIN cafes 
	ON products.cafe_id = cafes.id
	WHERE 1=1
	`

	args := []any{}
	argID := 1

	if productName != "" {
		query += fmt.Sprintf(" AND products.name ILIKE $%d", argID)
		args = append(args, "%"+productName+"%")
		argID++
	}

	query += " ORDER BY inventory.id"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argID)
		args = append(args, limit)
	}

	rows, err := i.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventories := []model.InventoryWithProduct{}

	for rows.Next() {
		var inventory model.InventoryWithProduct

		err := rows.Scan(
			&inventory.ID,
			&inventory.ProductId,
			&inventory.ProductName,
			&inventory.CafeId,
			&inventory.CafeName,
			&inventory.StockQuantity,
			&inventory.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		inventories = append(inventories, inventory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return inventories, nil
}

func (i *inventoryRepository) Create(ctx context.Context, inventory *model.Inventory) error {

	return i.db.QueryRow(
		ctx,
		`INSERT INTO inventory(
			product_id,
			stock_quantity
		)
		VALUES ($1, $2)
		RETURNING id
		`,
		inventory.ProductId,
		inventory.StockQuantity,
	).Scan(&inventory.ID)
}

func (i *inventoryRepository) UpdateStock(ctx context.Context, id int, changeValue int, isAdjust bool) (int, int, error) {

	var currentQty int
	err := i.db.QueryRow(ctx, `SELECT stock_quantity FROM inventory WHERE product_id = $1 FOR UPDATE`, id).Scan(&currentQty)
	if err != nil {
		return 0, 0, err
	}

	newQty := changeValue
	if !isAdjust {
		newQty = currentQty + changeValue
		if newQty < 0 {
			return 0, 0, ErrInsufficientStock
		}
	}

	var inventoryID, stockQuantity int
	err = i.db.QueryRow(ctx, `UPDATE inventory SET stock_quantity = $1 WHERE product_id = $2 RETURNING id, stock_quantity`, newQty, id).Scan(&inventoryID, &stockQuantity)

	return inventoryID, stockQuantity, err
}

func (i *inventoryRepository) Delete(ctx context.Context, id int) error {

	result, err := i.db.Exec(
		ctx,
		`DELETE FROM inventory
		WHERE product_id = $1`,
		id,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}
