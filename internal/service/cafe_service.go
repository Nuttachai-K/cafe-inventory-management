package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
	"github.com/jackc/pgx/v5"
)

var (
	ErrCafeNotFound = errors.New("cafe not found")
	ErrInvalidInput = errors.New("invalid input")
)

type CafeService interface {
	GetByID(ctx context.Context, id int) (*model.Cafe, error)
	GetAll(ctx context.Context, station string, limit int) ([]model.Cafe, error)
	Create(ctx context.Context, cafe *model.Cafe) error
	Update(ctx context.Context, id int, cu *model.CafeUpdate) error
	Delete(ctx context.Context, id int) error
}

type cafeService struct {
	repo repository.CafeRepository
}

func NewCafeService(repo repository.CafeRepository) CafeService {
	return &cafeService{
		repo: repo,
	}
}

func (s *cafeService) GetByID(ctx context.Context, id int) (*model.Cafe, error) {
	if id < 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	cafe, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, translateErr(err)
	}
	return cafe, nil
}

func (s *cafeService) GetAll(ctx context.Context, station string, limit int) ([]model.Cafe, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetAll(ctx, station, limit)
}

func (s *cafeService) Create(ctx context.Context, cafe *model.Cafe) error {
	if strings.TrimSpace(cafe.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if cafe.Latitude < -90 || cafe.Latitude > 90 {
		return fmt.Errorf("%w: invalid latitude", ErrInvalidInput)
	}
	if cafe.Longitude < -180 || cafe.Longitude > 180 {
		return fmt.Errorf("%w: invalid longitude", ErrInvalidInput)
	}
	if err := s.repo.Create(ctx, cafe); err != nil {
		return translateErr(err)
	}
	return nil
}

func (s *cafeService) Update(ctx context.Context, id int, cu *model.CafeUpdate) error {
	if id <= 0 {
		return fmt.Errorf("%w: invalid id", ErrInvalidInput)
	}
	if cu.Latitude != nil && (*cu.Latitude < -90 || *cu.Latitude > 90) {
		return fmt.Errorf("%w: invalid latitude", ErrInvalidInput)
	}
	if cu.Longitude != nil && (*cu.Longitude < -180 || *cu.Longitude > 180) {
		return fmt.Errorf("%w: invalid longitude", ErrInvalidInput)
	}
	if err := s.repo.Update(ctx, id, cu); err != nil {
		return translateErr(err)
	}
	return nil
}

func (s *cafeService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return translateErr(err)
	}
	return nil
}

func translateErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrCafeNotFound
	}
	return err
}
