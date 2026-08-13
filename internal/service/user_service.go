package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var ()

type UserService interface {
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetAll(ctx context.Context, limit int) ([]model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, id int, uu *model.UserUpdate) (*model.User, error)
	Delete(ctx context.Context, id int) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetByID(ctx context.Context, id int) (*model.User, error) {
	if id < 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: user", translateErr(err))
	}
	return user, nil
}

func (s *userService) GetAll(ctx context.Context, limit int) ([]model.User, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetAll(ctx, limit)
}

func (s *userService) Create(ctx context.Context, user *model.User) error {
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return fmt.Errorf("%w: email is not in correct format", ErrInvalidInput)
	}

	if !user.UserRole.Valid() {
		return ErrInvalidUserRole
	}

	if user.Password == "" {
		return fmt.Errorf("%w: password is required", ErrInvalidInput)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	user.PasswordHash = string(hash)
	user.Password = ""

	err = s.repo.Create(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		// Error from unique constraint in user email"
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: %s", ErrDuplicateEmail, user.Email)
		}
		return fmt.Errorf("%w: user", translateErr(err))
	}
	return nil
}

func (s *userService) Update(ctx context.Context, id int, uu *model.UserUpdate) (*model.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid id", ErrInvalidInput)
	}

	if uu.Email != nil {
		if _, err := mail.ParseAddress(*uu.Email); err != nil {
			return nil, fmt.Errorf("%w: email is not in correct format", ErrInvalidInput)
		}
	}

	if uu.UserRole != nil && !uu.UserRole.Valid() {
		return nil, ErrInvalidUserRole
	}

	if uu.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		hashed := string(hash)
		uu.PasswordHash = &hashed
		uu.Password = nil
	}

	user, err := s.repo.Update(ctx, id, uu)
	if err != nil {
		var pgErr *pgconn.PgError
		// Error from unique constraint in user email"
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateEmail, *uu.Email)
		}
		return nil, fmt.Errorf("%w: user", translateErr(err))
	}
	return user, nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return translateErr(err)
	}
	return nil
}
