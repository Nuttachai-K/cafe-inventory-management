package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type AuthService interface {
	Login(ctx context.Context, auth *model.Authentication) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Login(ctx context.Context, auth *model.Authentication) (string, error) {
	user, err := s.repo.GetByEmail(ctx, auth.Email)
	if err != nil {
		return "", fmt.Errorf("%w", ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(auth.Password)); err != nil {
		return "", ErrInvalidCredentials
	}
	return utils.GenerateToken(user.ID, string(user.UserRole))
}
