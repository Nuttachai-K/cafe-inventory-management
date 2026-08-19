package service

import (
	"context"
	"fmt"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryService interface {
	GetByID(ctx context.Context, id int) (*model.InventoryWithProduct, error)
	GetAll(ctx context.Context, productName string, limit int) ([]model.InventoryWithProduct, error)
	UpdateStock(ctx context.Context, id int, log *model.InventoryLog) (int, error)
}

type inventoryService struct {
	repo repository.InventoryRepository
	pool *pgxpool.Pool
}

func NewInventoryService(repo repository.InventoryRepository, pool *pgxpool.Pool) InventoryService {

	return &inventoryService{
		repo: repo,
		pool: pool,
	}
}

func (s *inventoryService) GetByID(ctx context.Context, id int) (*model.InventoryWithProduct, error) {

	if id < 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	inventory, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: inventory", translateErr(err))
	}
	return inventory, nil
}

func (s *inventoryService) GetAll(ctx context.Context, productName string, limit int) ([]model.InventoryWithProduct, error) {

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	inventories, err := s.repo.GetAll(ctx, productName, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: inventories", translateErr(err))
	}

	return inventories, nil
}

func (s *inventoryService) UpdateStock(ctx context.Context, id int, log *model.InventoryLog) (int, error) {

	if id < 0 {
		return 0, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	if log.ChangeQuantity <= 0 {
		return 0, fmt.Errorf("%w: change quantity must be positive", ErrInvalidInput)
	}

	op, err := model.ParseOperation(log.Operation)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}
	log.Operation = string(op)

	changeValue, isAdjust, err := resolveStockChange(op, log.ChangeQuantity)
	if err != nil {
		return 0, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	inventoryRepoTx := repository.NewInventoryRepository(tx)
	logRepoTx := repository.NewInventoryLogRepository(tx)

	inventoryID, stockQuantity, err := inventoryRepoTx.UpdateStock(ctx, id, changeValue, isAdjust)
	if err != nil {
		return 0, translateErr(err)
	}

	log.InventoryId = inventoryID

	if err = logRepoTx.Create(ctx, log); err != nil {
		return 0, translateErr(err)
	}

	return stockQuantity, tx.Commit(ctx)
}

func resolveStockChange(op model.Operation, quantity int) (changeValue int, isAdjust bool, err error) {

	switch op {
	case model.InOperation:
		return quantity, false, nil
	case model.OutOperation:
		return -quantity, false, nil
	case model.AdjustOperation:
		return quantity, true, nil
	default:
		return 0, false, fmt.Errorf("%w: operation must be IN, OUT, or ADJUST", ErrInvalidInput)
	}
}
