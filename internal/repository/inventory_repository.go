package repository

import (
	"context"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/jackc/pgx/v5"
)

type InventoryRepository interface {
	Create(ctx context.Context, inventory *model.Inventory) error
	Update(ctx context.Context, id int, quantity int) (*model.Inventory, error)
	Delete(ctx context.Context, id int) error
}

type inventoryRepository struct {
	db DBTX
}

func NewInventoryRepository(db DBTX) InventoryRepository {
	return &inventoryRepository{
		db: db,
	}
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

func (i *inventoryRepository) Update(ctx context.Context, id int, quantity int) (*model.Inventory, error) {

	result, err := i.db.Exec(
		ctx,
		"UPDATE inventory SET stock_quantity = $1 WHERE id = $2",
		quantity, id,
	)

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	var inventory model.Inventory
	err = i.db.QueryRow(
		ctx,
		`SELECT
			id,
			product_id,
			created_at,
			stock_quantity,
			updated_at
		FROM inventory
		WHERE id = $1`,
		id,
	).Scan(
		&inventory.ID,
		&inventory.ProductId,
		&inventory.StockQuantity,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &inventory, nil

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
