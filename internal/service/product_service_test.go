package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type mockProductRepo struct {
	getByIDFn func(ctx context.Context, id int) (*model.ProductWithCategory, error)
	getAllFn  func(ctx context.Context, limit int) ([]model.ProductWithCategory, error)
	createFn  func(ctx context.Context, user *model.Product) error
	updateFn  func(ctx context.Context, id int, pu *model.ProductUpdate) (*model.ProductWithCategory, error)
	deleteFn  func(ctx context.Context, id int) error
}

func (m *mockProductRepo) GetByID(ctx context.Context, id int) (*model.ProductWithCategory, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockProductRepo) GetAll(ctx context.Context, limit int) ([]model.ProductWithCategory, error) {
	return m.getAllFn(ctx, limit)
}
func (m *mockProductRepo) Create(ctx context.Context, user *model.Product) error {
	return m.createFn(ctx, user)
}
func (m *mockProductRepo) Update(ctx context.Context, id int, pu *model.ProductUpdate) (*model.ProductWithCategory, error) {
	return m.updateFn(ctx, id, pu)
}
func (m *mockProductRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFn(ctx, id)
}

func TestProductService_GetById(t *testing.T) {

	tests := []struct {
		name      string
		id        int
		getByIDFn func(ctx context.Context, id int) (*model.ProductWithCategory, error)
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
			getByIDFn: func(ctx context.Context, id int) (*model.ProductWithCategory, error) {
				return nil, pgx.ErrNoRows
			},
			wantErr: ErrDataNotFound,
		},
		{
			name: "found",
			id:   1,
			getByIDFn: func(ctx context.Context, id int) (*model.ProductWithCategory, error) {
				return &model.ProductWithCategory{
					ID:           1,
					CafeId:       2,
					Name:         "Latte",
					Price:        decimal.NewFromFloat(3.50),
					CategoryName: "Beverages",
				}, nil
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductService(&mockProductRepo{getByIDFn: tt.getByIDFn}, nil)
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

func TestProductService_GetAll(t *testing.T) {

	tests := []struct {
		name      string
		input     int
		wantLimit int
	}{
		{name: "zero use default", input: 0, wantLimit: 20},
		{name: "negative use default", input: -5, wantLimit: 20},
		{name: "over max clamps to 100", input: 500, wantLimit: 100},
		{name: "in range passes through", input: 50, wantLimit: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotLimit int
			repo := &mockProductRepo{
				getAllFn: func(ctx context.Context, limit int) ([]model.ProductWithCategory, error) {
					gotLimit = limit
					return nil, nil
				},
			}
			s := NewProductService(repo, nil)

			if _, err := s.GetAll(context.Background(), tt.input); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotLimit != tt.wantLimit {
				t.Fatalf("repo recieved limit=%d, want %d", gotLimit, tt.wantLimit)
			}
		})
	}
}

func TestProductService_Create(t *testing.T) {

	tests := []struct {
		name     string
		product  *model.Product
		createFn func(ctx context.Context, product *model.Product) error
		wantErr  error
	}{
		{
			name:    "invalid cafe id",
			product: &model.Product{CafeId: -1},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "invalid category id",
			product: &model.Product{CategoryId: -1},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "invalid price",
			product: &model.Product{Price: decimal.NewFromInt(-1)},
			wantErr: ErrInvalidInput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockProductRepo{createFn: tt.createFn}
			s := NewProductService(repo, nil)

			err := s.Create(context.Background(), tt.product)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want error wrapping %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductService_Update_PartialFields(t *testing.T) {

	repo := &mockProductRepo{
		updateFn: func(ctx context.Context, id int, pu *model.ProductUpdate) (*model.ProductWithCategory, error) {
			return &model.ProductWithCategory{
				ID: id,
			}, nil
		},
	}
	s := NewProductService(repo, nil)
	_, err := s.Update(context.Background(), 1, &model.ProductUpdate{
		Name: ptr("newname"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProductService_Delete(t *testing.T) {

	tests := []struct {
		name     string
		id       int
		deleteFn func(ctx context.Context, id int) error
		wantErr  error
	}{
		{name: "invalid id", id: 0, wantErr: ErrInvalidInput},
		{name: "success", id: 1, deleteFn: func(ctx context.Context, id int) error { return nil }, wantErr: nil},
		{name: "not found", id: 1, deleteFn: func(ctx context.Context, id int) error { return pgx.ErrNoRows }, wantErr: ErrDataNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductService(&mockProductRepo{deleteFn: tt.deleteFn}, nil)
			err := s.Delete(context.Background(), tt.id)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("got %v, want error wrapping %v", err, tt.wantErr)
			}
		})
	}
}
