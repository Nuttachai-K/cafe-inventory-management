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
		return nil, fmt.Errorf("%w: cafe", translateErr(err))
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

	return s.repo.GetAll(ctx, productName, limit)
}

func (s *inventoryService) UpdateStock(ctx context.Context, id int, log *model.InventoryLog) (int, error) {

	if id < 0 {
		return 0, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	if log.ChangeQuantity <= 0 {
		return 0, fmt.Errorf("%w: change quantity must be positive", ErrInvalidInput)
	}

	var changeValue int
	isAdjust := false
	switch log.Operation {
	case "IN":
		changeValue = log.ChangeQuantity
	case "OUT":
		changeValue = -log.ChangeQuantity
	case "ADJUST":
		changeValue = log.ChangeQuantity
		isAdjust = true
	default:
		return 0, fmt.Errorf("%w: operation must be IN, OUT, or ADJUST", ErrInvalidInput)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	inventoryRepoTx := repository.NewInventoryRepository(tx)
	logRepoTx := repository.NewInventoryLogRepository(tx)

	stockQuantity, err := inventoryRepoTx.UpdateStock(ctx, id, changeValue, isAdjust)
	if err != nil {
		return 0, translateErr(err)
	}

	if err = logRepoTx.Create(ctx, log); err != nil {
		return 0, translateErr(err)
	}

	return stockQuantity, tx.Commit(ctx)
}
