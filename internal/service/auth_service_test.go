package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.GetByEmailFn(ctx, email)
}

func TestAuthService(t *testing.T) {
	tests := []struct {
		name         string
		auth         *model.Authentication
		GetByEmailFn func(ctx context.Context, email string) (*model.User, error)
		wantErr      error
	}{
		{
			name: "user not found",
			auth: &model.Authentication{
				Email:    "test",
				Password: "123",
			},
			GetByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return nil, pgx.ErrNoRows
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			auth: &model.Authentication{
				Email:    "test@example.com",
				Password: "wrong",
			},
			GetByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
				return &model.User{ID: 1, PasswordHash: string(hash), UserRole: model.RoleStaff, IsActive: true}, nil
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "user inactived",
			auth: &model.Authentication{
				Email:    "test@example.com",
				Password: "correct",
			},
			GetByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
				return &model.User{ID: 1, PasswordHash: string(hash), UserRole: model.RoleStaff, IsActive: false}, nil
			},
			wantErr: ErrUserInactive,
		},
		{
			name: "success",
			auth: &model.Authentication{
				Email:    "test@example.com",
				Password: "correct",
			},
			GetByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
				return &model.User{ID: 1, PasswordHash: string(hash), UserRole: model.RoleStaff, IsActive: true}, nil
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAuthService(&mockUserRepo{GetByEmailFn: tt.GetByEmailFn})
			token, err := s.Login(context.Background(), tt.auth)

			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("got %v, want error wrapping %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && token == "" {
				t.Fatal("expected non-empty token on success")
			}
		})
	}
}
