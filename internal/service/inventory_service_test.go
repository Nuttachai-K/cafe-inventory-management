package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/jackc/pgx/v5"
)

type mockInventoryRepo struct {
	getByIDFn func(ctx context.Context, id int) (*model.InventoryWithProduct, error)
	getAllFn  func(ctx context.Context, productName string, limit int) ([]model.InventoryWithProduct, error)
	createFn  func(ctx context.Context, inventory *model.Inventory) error
	updateFn  func(ctx context.Context, id int, quantity int, isAdjust bool) (int, error)
	deleteFn  func(ctx context.Context, id int) error
}

func (m *mockInventoryRepo) GetByID(ctx context.Context, id int) (*model.InventoryWithProduct, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockInventoryRepo) GetAll(ctx context.Context, productName string, limit int) ([]model.InventoryWithProduct, error) {
	return m.getAllFn(ctx, productName, limit)
}
func (m *mockInventoryRepo) Create(ctx context.Context, inventory *model.Inventory) error {
	return m.createFn(ctx, inventory)
}
func (m *mockInventoryRepo) UpdateStock(ctx context.Context, id int, quantity int, isAdjust bool) (int, error) {
	return m.updateFn(ctx, id, quantity, isAdjust)
}
func (m *mockInventoryRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFn(ctx, id)
}

type mockInventoryLogRepo struct {
	getAllFn func(ctx context.Context, limit int) ([]model.InventoryLog, error)
	createFn func(ctx context.Context, inventoryLog *model.InventoryLog) error
}

func (ml *mockInventoryLogRepo) GetAll(ctx context.Context, limit int) ([]model.InventoryLog, error) {
	return ml.getAllFn(ctx, limit)
}
func (ml *mockInventoryLogRepo) Create(ctx context.Context, inventoryLog *model.InventoryLog) error {
	return ml.createFn(ctx, inventoryLog)
}

func TestInventoryService_GetById(t *testing.T) {

	tests := []struct {
		name      string
		id        int
		getByIDFn func(ctx context.Context, id int) (*model.InventoryWithProduct, error)
		wantErr   error
	}{
		{
			name:    "negaive id",
			id:      -1,
			wantErr: ErrInvalidInput,
		},
		{
			name: "not found",
			id:   99,
			getByIDFn: func(ctx context.Context, id int) (*model.InventoryWithProduct, error) {
				return nil, pgx.ErrNoRows
			},
			wantErr: ErrDataNotFound,
		},
		{
			name: "found",
			id:   1,
			getByIDFn: func(ctx context.Context, id int) (*model.InventoryWithProduct, error) {
				return &model.InventoryWithProduct{
					ID:            1,
					ProductId:     1,
					ProductName:   "Matcha cake",
					CafeId:        2,
					CafeName:      "Tokyo cafe",
					StockQuantity: 10,
				}, nil
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInventoryService(&mockInventoryRepo{getByIDFn: tt.getByIDFn}, nil)
			_, err := s.GetByID(context.Background(), tt.id)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("got %v, want error wrapping %v", err, tt.wantErr)
			}
		})
	}
}

func TestInventoryService_GetAll(t *testing.T) {

	tests := []struct {
		name        string
		productName string
		input       int
		wantLimit   int
	}{
		{name: "zero use default", productName: "Test 1", input: 0, wantLimit: 20},
		{name: "negative use default", productName: "Test 2", input: -5, wantLimit: 20},
		{name: "over max clamps to 100", productName: "Test 3", input: 500, wantLimit: 100},
		{name: "in range passes through", productName: "Test 4", input: 50, wantLimit: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotLimit int
			repo := &mockInventoryRepo{
				getAllFn: func(ctx context.Context, productName string, limit int) ([]model.InventoryWithProduct, error) {
					gotLimit = limit
					return nil, nil
				},
			}
			s := NewInventoryService(repo, nil)

			if _, err := s.GetAll(context.Background(), tt.productName, tt.input); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotLimit != tt.wantLimit {
				t.Fatalf("repo recieved limit=%d, want %d", gotLimit, tt.wantLimit)
			}
		})
	}
}

func TestInventoryService_Update_Stock(t *testing.T) {

	tests := []struct {
		name    string
		id      int
		log     *model.InventoryLog
		wantErr error
	}{
		{name: "negative id", id: -1, log: &model.InventoryLog{ChangeQuantity: 5}, wantErr: ErrInvalidInput},
		{name: "zero change quantity", id: 1, log: &model.InventoryLog{ChangeQuantity: 0}, wantErr: ErrInvalidInput},
		{name: "negative change quantity", id: 1, log: &model.InventoryLog{ChangeQuantity: -3}, wantErr: ErrInvalidInput},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewInventoryService(&mockInventoryRepo{}, nil)
			_, err := s.UpdateStock(context.Background(), tt.id, tt.log)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got %v, want error wrapping %v", err, tt.wantErr)
			}
		})
	}
}
