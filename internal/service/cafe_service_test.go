package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockCafeRepo struct {
	getByIDFn func(ctx context.Context, id int) (*model.Cafe, error)
	getAllFn  func(ctx context.Context, filter model.CafeFilter) ([]model.Cafe, error)
	createFn  func(ctx context.Context, cafe *model.Cafe) error
	updateFn  func(ctx context.Context, id int, cu *model.CafeUpdate) (*model.Cafe, error)
	deleteFn  func(ctx context.Context, id int) error
}

func (m *mockCafeRepo) GetByID(ctx context.Context, id int) (*model.Cafe, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockCafeRepo) GetAll(ctx context.Context, filter model.CafeFilter) ([]model.Cafe, error) {
	return m.getAllFn(ctx, filter)
}

func (m *mockCafeRepo) Create(ctx context.Context, cafe *model.Cafe) error {
	return m.createFn(ctx, cafe)
}

func (m *mockCafeRepo) Update(ctx context.Context, id int, cu *model.CafeUpdate) (*model.Cafe, error) {
	return m.updateFn(ctx, id, cu)
}

func (m *mockCafeRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFn(ctx, id)
}

func TestCafeService_GetById(t *testing.T) {

	tests := []struct {
		name      string
		id        int
		getByIDFn func(ctx context.Context, id int) (*model.Cafe, error)
		wantErr   error
	}{
		{
			name:    "negative id",
			id:      -1,
			wantErr: ErrInvalidInput,
		},
		{
			name: "not found",
			id:   99,
			getByIDFn: func(ctx context.Context, id int) (*model.Cafe, error) {
				return nil, pgx.ErrNoRows
			},
			wantErr: ErrDataNotFound,
		},
		{
			name: "found",
			id:   1,
			getByIDFn: func(ctx context.Context, id int) (*model.Cafe, error) {
				return &model.Cafe{ID: 1, Name: "Blue bee"}, nil
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewCafeService(&mockCafeRepo{getByIDFn: tt.getByIDFn})
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

func TestCafeService_GetAll_Limit(t *testing.T) {

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
			repo := &mockCafeRepo{
				getAllFn: func(ctx context.Context, filter model.CafeFilter) ([]model.Cafe, error) {
					gotLimit = filter.Limit
					return nil, nil
				},
			}
			s := NewCafeService(repo)

			_, err := s.GetAll(context.Background(), model.CafeFilter{Limit: tt.input})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotLimit != tt.wantLimit {
				t.Fatalf("repo received limit=%d, want %d", gotLimit, tt.wantLimit)
			}
		})
	}
}

func TestCafeService_GetAll_GeoValidation(t *testing.T) {

	tests := []struct {
		name    string
		filter  model.CafeFilter
		wantErr error
	}{
		{
			name:    "lat without lng",
			filter:  model.CafeFilter{Lat: new(35.0)},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "lng without lat",
			filter:  model.CafeFilter{Lng: new(139.0)},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "latitude out of range",
			filter:  model.CafeFilter{Lat: new(200.0), Lng: new(139.0)},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "longitude out of range",
			filter:  model.CafeFilter{Lat: new(35.0), Lng: new(200.0)},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "radius without lat/lng",
			filter:  model.CafeFilter{RadiusKm: new(5.0)},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "radius not positive",
			filter:  model.CafeFilter{RadiusKm: new(-10.0)},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "valid lat/lng/radius",
			filter:  model.CafeFilter{Lat: new(35.0), Lng: new(139.0), RadiusKm: new(5.0)},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCafeRepo{
				getAllFn: func(ctx context.Context, filter model.CafeFilter) ([]model.Cafe, error) {
					return nil, nil
				},
			}

			s := NewCafeService(repo)

			_, err := s.GetAll(context.Background(), tt.filter)
			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("got %v, want error wrapping %v", err, tt.wantErr)
			}
		})
	}
}

func TestCafeService_Create(t *testing.T) {

	tests := []struct {
		name    string
		cafe    *model.Cafe
		wantErr error
	}{
		{
			name:    "empty name",
			cafe:    &model.Cafe{Name: " ", Latitude: 35, Longitude: 139},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "invalid latitude",
			cafe:    &model.Cafe{Name: "Cafe A", Latitude: 200, Longitude: 139},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "invalid longitude",
			cafe:    &model.Cafe{Name: "Cafe A", Latitude: 35, Longitude: 200},
			wantErr: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCafeRepo{
				createFn: func(ctx context.Context, cafe *model.Cafe) error {
					return nil
				},
			}
			s := NewCafeService(repo)

			err := s.Create(context.Background(), tt.cafe)
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

func TestCafeService_Update(t *testing.T) {

	tests := []struct {
		name    string
		id      int
		update  *model.CafeUpdate
		wantErr error
	}{
		{
			name:    "invalid id",
			id:      0,
			update:  &model.CafeUpdate{},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "invalid latitude",
			id:      1,
			update:  &model.CafeUpdate{Latitude: new(200.0)},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "invalid longitude",
			id:      1,
			update:  &model.CafeUpdate{Longitude: new(200.0)},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "success",
			id:      1,
			update:  &model.CafeUpdate{Name: new("New name")},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCafeRepo{
				updateFn: func(ctx context.Context, id int, cu *model.CafeUpdate) (*model.Cafe, error) {
					return &model.Cafe{ID: id}, nil
				},
			}

			s := NewCafeService(repo)

			_, err := s.Update(context.Background(), tt.id, tt.update)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("got %v, want error wrapping %v", err, tt.wantErr)
			}
		})
	}
}

func TestCafeService_Delete(t *testing.T) {

	tests := []struct {
		name     string
		id       int
		deleteFn func(ctx context.Context, id int) error
		wantErr  error
	}{
		{
			name:    "invalid id",
			id:      0,
			wantErr: ErrInvalidInput,
		},
		{
			name: "has dependent products",
			id:   1,
			deleteFn: func(ctx context.Context, id int) error {
				return &pgconn.PgError{Code: "23503"}
			},
			wantErr: ErrHasDependents,
		},
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
			repo := &mockCafeRepo{
				deleteFn: tt.deleteFn,
			}
			s := NewCafeService(repo)

			err := s.Delete(context.Background(), tt.id)
			if err != nil && tt.wantErr == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("got: %v, want error wrapping %v", err, tt.wantErr)
			}
		})
	}
}
