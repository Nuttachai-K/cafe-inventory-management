package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockCategoryRepo struct {
	getByIDFn func(ctx context.Context, id int) (*model.Category, error)
	getAllFn  func(ctx context.Context, limit int) ([]model.Category, error)
	createFn  func(ctx context.Context, category *model.Category) error
	updateFn  func(ctx context.Context, id int, name string) (*model.Category, error)
	deleteFn  func(ctx context.Context, id int) error
}

func (m *mockCategoryRepo) GetByID(ctx context.Context, id int) (*model.Category, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockCategoryRepo) GetAll(ctx context.Context, limit int) ([]model.Category, error) {
	return m.getAllFn(ctx, limit)
}
func (m *mockCategoryRepo) Create(ctx context.Context, category *model.Category) error {
	return m.createFn(ctx, category)
}
func (m *mockCategoryRepo) Update(ctx context.Context, id int, name string) (*model.Category, error) {
	return m.updateFn(ctx, id, name)
}
func (m *mockCategoryRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFn(ctx, id)
}

func TestCategoryService_GetById(t *testing.T) {

	tests := []struct {
		name      string
		id        int
		getByIDFn func(ctx context.Context, id int) (*model.Category, error)
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
			getByIDFn: func(ctx context.Context, id int) (*model.Category, error) {
				return nil, pgx.ErrNoRows
			},
			wantErr: ErrDataNotFound,
		},
		{
			name: "found",
			id:   1,
			getByIDFn: func(ctx context.Context, id int) (*model.Category, error) {
				return &model.Category{
					ID:   1,
					Name: "Latte",
				}, nil
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewCategoryService(&mockCategoryRepo{getByIDFn: tt.getByIDFn})
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

func TestCategoryService_GetAll(t *testing.T) {

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
			repo := &mockCategoryRepo{
				getAllFn: func(ctx context.Context, limit int) ([]model.Category, error) {
					gotLimit = limit
					return nil, nil
				},
			}
			s := NewCategoryService(repo)

			if _, err := s.GetAll(context.Background(), tt.input); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotLimit != tt.wantLimit {
				t.Fatalf("repo recieved limit=%d, want %d", gotLimit, tt.wantLimit)
			}
		})
	}
}

func TestCategoryService_Create(t *testing.T) {

	tests := []struct {
		name     string
		category *model.Category
		createFn func(ctx context.Context, product *model.Category) error
		wantErr  error
	}{
		{
			name:     "Duplicate category's name",
			category: &model.Category{Name: "a"},
			createFn: func(ctx context.Context, category *model.Category) error {
				return &pgconn.PgError{Code: "23505"}
			},
			wantErr: ErrDuplicateCategory,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCategoryRepo{createFn: tt.createFn}
			s := NewCategoryService(repo)

			err := s.Create(context.Background(), tt.category)

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

func TestCategoryService_Update(t *testing.T) {

	tests := []struct {
		name         string
		id           int
		categoryName string
		wantErr      error
	}{
		{
			name:    "invalid id",
			id:      -4,
			wantErr: ErrInvalidInput,
		},
		{
			name:         "success",
			id:           1,
			categoryName: "New category",
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCategoryRepo{
				updateFn: func(ctx context.Context, id int, name string) (*model.Category, error) {
					return &model.Category{
						ID: id,
					}, nil
				},
			}
			s := NewCategoryService(repo)
			_, err := s.Update(context.Background(), tt.id, tt.categoryName)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("got: %v, want error wrapping: %v", err, tt.wantErr)
			}
		})
	}

}

func TestCategoryService_Delete(t *testing.T) {

	tests := []struct {
		name     string
		id       int
		deleteFn func(ctx context.Context, id int) error
		wantErr  error
	}{
		{name: "invalid id", id: 0, wantErr: ErrInvalidInput},
		{
			name: "not found",
			id:   99,
			deleteFn: func(ctx context.Context, id int) error {
				return pgx.ErrNoRows
			},
			wantErr: ErrDataNotFound,
		},
		{
			name: "success",
			id:   1,
			deleteFn: func(ctx context.Context, id int) error {
				return nil
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewCategoryService(&mockCategoryRepo{deleteFn: tt.deleteFn})
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
