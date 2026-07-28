package service

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrDataNotFound    = errors.New("data not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrDuplicateEmail  = errors.New("This email is already used")
	ErrInvalidUserRole = errors.New("invalid user role")
)

func translateErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrDataNotFound
	}
	return err
}
