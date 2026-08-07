package service

import (
	"context"
	"fmt"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
)

type InventoryLogService interface {
	GetLog(ctx context.Context, filter *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error)
}

type inventoryLogService struct {
	repo repository.InventoryLogRepository
}

func NewInventoryLogService(repo repository.InventoryLogRepository) InventoryLogService {
	return &inventoryLogService{
		repo: repo,
	}
}

func (s *inventoryLogService) GetLog(ctx context.Context, filter *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error) {

	if filter.InventoryId != nil && *filter.InventoryId < 0 {
		return nil, fmt.Errorf("%w: inventory id must be positive", ErrInvalidInput)
	}

	if filter.UserId != nil && *filter.UserId < 0 {
		return nil, fmt.Errorf("%w: user id must be positive", ErrInvalidInput)
	}

	if filter.Operation != nil {
		op, err := model.ParseOperation(*filter.Operation)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
		*filter.Operation = string(op)
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	if filter.From != nil && filter.To != nil && filter.From.After(*filter.To) {
		return nil, fmt.Errorf("%w: from must be before to", ErrInvalidInput)
	}

	inventoryLogs, err := s.repo.GetLog(ctx, filter, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: logs", translateErr(err))
	}

	return inventoryLogs, nil
}
