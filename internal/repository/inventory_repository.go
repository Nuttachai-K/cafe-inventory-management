package repository

import (
	"context"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
)

type InventoryRepository interface {
	Create(ctx context.Context, inventory *model.Inventory) error
	// Delete(ctx context.Context, id int) error
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
