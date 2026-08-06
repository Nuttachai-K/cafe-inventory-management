package service

import (
	"context"
	"fmt"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
)

type InventoryLogService interface {
	GetLog(ctx context.Context, filters *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error)
}

type inventoryLogService struct {
	repo repository.InventoryLogRepository
}

func (s *inventoryLogService) GetLog(ctx context.Context, filters *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error) {

	if filters.InventoryId != nil && *filters.InventoryId < 0 {
		return nil, fmt.Errorf("%w: inventory id must be positive", ErrInvalidInput)
	}

	if filters.UserId != nil && *filters.UserId < 0 {
		return nil, fmt.Errorf("%w: user id must be positive", ErrInvalidInput)
	}

	if filters.Operation != nil {
		op, err := model.ParseOperation(*filters.Operation)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
		*filters.Operation = string(op)
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	if filters.From != nil && filters.To != nil && filters.From.After(*filters.To) {
		return nil, fmt.Errorf("%w: from must be before to", ErrInvalidInput)
	}

	inventoryLogs, err := s.repo.GetAll(ctx, filters, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: logs", translateErr(err))
	}

	return inventoryLogs, nil
}
