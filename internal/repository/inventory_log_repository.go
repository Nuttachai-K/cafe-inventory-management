package repository

import (
	"context"
	"fmt"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
)

type InventoryLogRepository interface {
	GetLog(ctx context.Context, log *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error)
	Create(ctx context.Context, inventoryLog *model.InventoryLog) error
}

type inventoryLogRepository struct {
	db DBTX
}

func NewInventoryLogRepository(db DBTX) InventoryLogRepository {

	return &inventoryLogRepository{
		db: db,
	}
}

func (il *inventoryLogRepository) GetLog(ctx context.Context, filter *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error) {

	query := `
	SELECT
		id,
		inventory_id,
		user_id,
		change_quantity,
		operation,
		created_at
	FROM inventory_logs
	WHERE 1=1
	`

	args := []any{}
	argID := 1

	if filter.InventoryId != nil {
		query += fmt.Sprintf(" AND inventory_id = $%d", argID)
		args = append(args, *filter.InventoryId)
		argID++
	}

	if filter.UserId != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argID)
		args = append(args, *filter.UserId)
		argID++
	}

	if filter.Operation != nil {
		query += fmt.Sprintf(" AND operation = $%d", argID)
		args = append(args, *filter.Operation)
		argID++
	}

	if filter.From != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argID)
		args = append(args, *filter.From)
		argID++
	}

	if filter.To != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argID)
		args = append(args, *filter.To)
		argID++
	}

	query += " ORDER BY id"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argID)
		args = append(args, limit)
	}

	rows, err := il.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventoryLogs := []model.InventoryLog{}

	for rows.Next() {
		var inventoryLog model.InventoryLog

		err := rows.Scan(
			&inventoryLog.ID,
			&inventoryLog.InventoryId,
			&inventoryLog.UserId,
			&inventoryLog.ChangeQuantity,
			&inventoryLog.Operation,
			&inventoryLog.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		inventoryLogs = append(inventoryLogs, inventoryLog)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return inventoryLogs, nil
}

func (il *inventoryLogRepository) Create(ctx context.Context, inventoryLog *model.InventoryLog) error {

	return il.db.QueryRow(
		ctx,
		`INSERT INTO inventory_logs (
			inventory_id,
			user_id,
			change_quantity,
			operation
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`,
		inventoryLog.InventoryId,
		inventoryLog.UserId,
		inventoryLog.ChangeQuantity,
		inventoryLog.Operation,
	).Scan(&inventoryLog.ID)
}
