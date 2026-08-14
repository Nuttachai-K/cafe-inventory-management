package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
)

type CafeService interface {
	GetByID(ctx context.Context, id int) (*model.Cafe, error)
	GetAll(ctx context.Context, filter model.CafeFilter) ([]model.Cafe, error)
	Create(ctx context.Context, cafe *model.Cafe) error
	Update(ctx context.Context, id int, cu *model.CafeUpdate) (*model.Cafe, error)
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
		return nil, fmt.Errorf("%w: cafe", translateErr(err))
	}
	return cafe, nil
}

func (s *cafeService) GetAll(ctx context.Context, filter model.CafeFilter) ([]model.Cafe, error) {

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if (filter.Lat != nil) != (filter.Lng != nil) {
		return nil, fmt.Errorf("%w: lat and lng must be provided together", ErrInvalidInput)
	}
	if filter.Lat != nil && (*filter.Lat < -90 || *filter.Lat > 90) {
		return nil, fmt.Errorf("%w: invalid latitude", ErrInvalidInput)
	}
	if filter.Lng != nil && (*filter.Lng < -180 || *filter.Lng > 180) {
		return nil, fmt.Errorf("%w: invalid longitude", ErrInvalidInput)
	}
	if filter.RadiusKm != nil {
		if filter.Lat == nil {
			return nil, fmt.Errorf("%w: radius requires lat and lng", ErrInvalidInput)
		}
		if *filter.RadiusKm <= 0 {
			return nil, fmt.Errorf("%w: radius must be positive", ErrInvalidInput)
		}
	}

	return s.repo.GetAll(ctx, filter)
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
		return fmt.Errorf("%w: cafe", translateErr(err))
	}
	return nil
}

func (s *cafeService) Update(ctx context.Context, id int, cu *model.CafeUpdate) (*model.Cafe, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid id", ErrInvalidInput)
	}
	if cu.Latitude != nil && (*cu.Latitude < -90 || *cu.Latitude > 90) {
		return nil, fmt.Errorf("%w: invalid latitude", ErrInvalidInput)
	}
	if cu.Longitude != nil && (*cu.Longitude < -180 || *cu.Longitude > 180) {
		return nil, fmt.Errorf("%w: invalid longitude", ErrInvalidInput)
	}
	cafe, err := s.repo.Update(ctx, id, cu)
	if err != nil {
		return nil, fmt.Errorf("%w: cafe", translateErr(err))
	}
	return cafe, nil
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
