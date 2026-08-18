package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockUserRepo struct {
	getByIDFn    func(ctx context.Context, id int) (*model.User, error)
	getAllFn     func(ctx context.Context, limit int) ([]model.User, error)
	GetByEmailFn func(ctx context.Context, email string) (*model.User, error)
	createFn     func(ctx context.Context, user *model.User) error
	updateFn     func(ctx context.Context, id int, uu *model.UserUpdate) (*model.User, error)
	deleteFn     func(ctx context.Context, id int) error
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockUserRepo) GetAll(ctx context.Context, limit int) ([]model.User, error) {
	return m.getAllFn(ctx, limit)
}
func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error {
	return m.createFn(ctx, user)
}
func (m *mockUserRepo) Update(ctx context.Context, id int, uu *model.UserUpdate) (*model.User, error) {
	return m.updateFn(ctx, id, uu)
}
func (m *mockUserRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFn(ctx, id)
}

func TestUserService_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		getByIDFn func(ctx context.Context, id int) (*model.User, error)
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
			getByIDFn: func(ctx context.Context, id int) (*model.User, error) {
				return nil, pgx.ErrNoRows
			},
			wantErr: ErrDataNotFound,
		},
		{
			name: "found",
			id:   1,
			getByIDFn: func(ctx context.Context, id int) (*model.User, error) {
				return &model.User{ID: id}, nil
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(&mockUserRepo{getByIDFn: tt.getByIDFn})
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

func TestUserService_GetAll(t *testing.T) {
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
			repo := &mockUserRepo{
				getAllFn: func(ctx context.Context, limit int) ([]model.User, error) {
					gotLimit = limit
					return nil, nil
				},
			}
			s := NewUserService(repo)

			if _, err := s.GetAll(context.Background(), tt.input); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotLimit != tt.wantLimit {
				t.Fatalf("repo received limit=%d, want %d", gotLimit, tt.wantLimit)
			}
		})
	}
}

func TestUserService_Create(t *testing.T) {
	tests := []struct {
		name     string
		user     *model.User
		createFn func(ctx context.Context, user *model.User) error
		wantErr  error // sentinel to check with errors.Is; nil means "no error expected"
	}{
		{
			name: "valid user",
			user: &model.User{Email: "a@example.com", UserRole: model.RoleStaff, Password: "a"},
			createFn: func(ctx context.Context, user *model.User) error {
				return nil
			},
			wantErr: nil,
		},
		{
			name:    "invalid email",
			user:    &model.User{Email: "not-an-email", UserRole: model.RoleStaff, Password: "a"},
			wantErr: ErrInvalidInput,
		},
		{
			name:    "invalid role",
			user:    &model.User{Email: "a@example.com", UserRole: "OWNER", Password: "a"},
			wantErr: ErrInvalidUserRole,
		},
		{
			name: "duplicate email",
			user: &model.User{Email: "a@example.com", UserRole: model.RoleStaff, Password: "a"},
			createFn: func(ctx context.Context, user *model.User) error {
				return &pgconn.PgError{Code: "23505"}
			},
			wantErr: ErrDuplicateEmail,
		},
		{
			name:    "Password is empty",
			user:    &model.User{Email: "b@example.com", UserRole: model.RoleStaff, Password: ""},
			wantErr: ErrInvalidInput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockUserRepo{createFn: tt.createFn}
			s := NewUserService(repo)

			err := s.Create(context.Background(), tt.user)

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

func TestUserService_Update_PartialFields(t *testing.T) {
	repo := &mockUserRepo{
		updateFn: func(ctx context.Context, id int, uu *model.UserUpdate) (*model.User, error) {
			return &model.User{ID: id}, nil
		},
	}
	s := NewUserService(repo)

	// Email and UserRole both nil — must not panic, must not error.
	_, err := s.Update(context.Background(), 1, &model.UserUpdate{
		Username: new("newname"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserService_Delete(t *testing.T) {
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
			s := NewUserService(&mockUserRepo{deleteFn: tt.deleteFn})
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
