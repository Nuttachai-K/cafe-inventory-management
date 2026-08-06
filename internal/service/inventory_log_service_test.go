package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
)

type mockInventoryLogRepo struct {
	getLogFn func(ctx context.Context, filter *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error)
	createFn func(ctx context.Context, inventoryLog *model.InventoryLog) error
}

func (m *mockInventoryLogRepo) GetLog(ctx context.Context, filter *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error) {
	return m.getLogFn(ctx, filter, limit)
}
func (m *mockInventoryLogRepo) Create(ctx context.Context, inventoryLog *model.InventoryLog) error {
	return m.createFn(ctx, inventoryLog)
}

func TestInventoryLogService_GetLog_Validation(t *testing.T) {
	tests := []struct {
		name   string
		filter *model.InventoryLogFilter
	}{
		{name: "negative inventory id", filter: &model.InventoryLogFilter{InventoryId: new(-1)}},
		{name: "negative user id", filter: &model.InventoryLogFilter{UserId: new(-1)}},
		{name: "invalid operation", filter: &model.InventoryLogFilter{Operation: new("test")}},
		{
			name: "from after to",
			filter: &model.InventoryLogFilter{
				From: new(time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)),
				To:   new(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockInventoryLogRepo{
				getLogFn: func(ctx context.Context, filter *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error) {
					t.Fatal("repo should not be called when validation fails")
					return nil, nil
				},
			}
			s := NewInventoryLogService(repo)

			_, err := s.GetLog(context.Background(), tt.filter, 20)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("got %v, want error wrapping %v", err, ErrInvalidInput)
			}
		})
	}
}

func TestInventoryLogService_GetLog_Limit(t *testing.T) {
	tests := []struct {
		name      string
		input     int
		wantLimit int
	}{
		{name: "zero uses default", input: 0, wantLimit: 20},
		{name: "negative uses default", input: -5, wantLimit: 20},
		{name: "over max clamps to 100", input: 500, wantLimit: 100},
		{name: "in range passes through", input: 50, wantLimit: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotLimit int
			repo := &mockInventoryLogRepo{
				getLogFn: func(ctx context.Context, filter *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error) {
					gotLimit = limit
					return nil, nil
				},
			}
			s := NewInventoryLogService(repo)

			if _, err := s.GetLog(context.Background(), &model.InventoryLogFilter{}, tt.input); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotLimit != tt.wantLimit {
				t.Fatalf("repo received limit=%d, want %d", gotLimit, tt.wantLimit)
			}
		})
	}
}

func TestInventoryLogService_GetLog_NormalizesOperation(t *testing.T) {
	var gotOperation *string
	repo := &mockInventoryLogRepo{
		getLogFn: func(ctx context.Context, filter *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error) {
			gotOperation = filter.Operation
			return nil, nil
		},
	}
	s := NewInventoryLogService(repo)

	filter := &model.InventoryLogFilter{Operation: new("in")}
	if _, err := s.GetLog(context.Background(), filter, 20); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotOperation == nil || *gotOperation != "IN" {
		t.Fatalf("got operation %v, want %q", gotOperation, "IN")
	}
}

func TestInventoryLogService_GetLog_ReturnsLogs(t *testing.T) {
	want := []model.InventoryLog{
		{ID: 1, InventoryId: 2, UserId: 3, ChangeQuantity: 5, Operation: "IN"},
	}
	repo := &mockInventoryLogRepo{
		getLogFn: func(ctx context.Context, filter *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error) {
			return want, nil
		},
	}
	s := NewInventoryLogService(repo)

	got, err := s.GetLog(context.Background(), &model.InventoryLogFilter{}, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != want[0].ID {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestInventoryLogService_GetLog_RepoError(t *testing.T) {
	dbErr := errors.New("db error")
	repo := &mockInventoryLogRepo{
		getLogFn: func(ctx context.Context, filter *model.InventoryLogFilter, limit int) ([]model.InventoryLog, error) {
			return nil, dbErr
		},
	}
	s := NewInventoryLogService(repo)

	_, err := s.GetLog(context.Background(), &model.InventoryLogFilter{}, 20)
	if !errors.Is(err, dbErr) {
		t.Fatalf("got %v, want error wrapping %v", err, dbErr)
	}
}
