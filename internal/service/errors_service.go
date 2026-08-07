package service

import (
	"errors"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
	"github.com/jackc/pgx/v5"
)

var (
	ErrDataNotFound       = errors.New("data not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrDuplicateEmail     = errors.New("this email is already used")
	ErrDuplicateCategory  = errors.New("this category name is already used")
	ErrInsufficientStock  = errors.New("the stock doesnt have enough items")
	ErrInvalidUserRole    = errors.New("invalid user role")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

func translateErr(err error) error {
	if errors.Is(err, repository.ErrInsufficientStock) {
		return ErrInsufficientStock
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrDataNotFound
	}
	return err
}
