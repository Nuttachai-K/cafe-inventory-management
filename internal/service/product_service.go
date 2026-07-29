package service

import (
	"context"
	"fmt"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
)

type ProductService interface {
	GetByID(ctx context.Context, id int) (*model.ProductWithCategory, error)
	GetAll(ctx context.Context, limit int) ([]model.ProductWithCategory, error)
	Create(ctx context.Context, user *model.Product) error
	Update(ctx context.Context, id int, pu *model.ProductUpdate) (*model.ProductWithCategory, error)
	Delete(ctx context.Context, id int) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetByID(ctx context.Context, id int) (*model.ProductWithCategory, error) {

	if id < 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: product", translateErr(err))
	}
	return product, nil
}

func (s *productService) GetAll(ctx context.Context, limit int) ([]model.ProductWithCategory, error) {

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetAll(ctx, limit)
}

func (s *productService) Create(ctx context.Context, product *model.Product) error {

	if product.CafeId < 0 {
		return fmt.Errorf("%w: user id must be positive", ErrInvalidInput)
	}
	if product.CategoryId < 0 {
		return fmt.Errorf("%w: user id must be positive", ErrInvalidInput)
	}
	if product.Price.IsNegative() {
		return fmt.Errorf("%w: price must be positive", ErrInvalidInput)
	}

	return s.repo.Create(ctx, product)
}

func (s *productService) Update(ctx context.Context, id int, pu *model.ProductUpdate) (*model.ProductWithCategory, error) {

	if id < 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	if pu.CafeId != nil {
		if *pu.CafeId < 0 {
			return nil, fmt.Errorf("%w: user id must be positive", ErrInvalidInput)
		}
	}
	if pu.CategoryId != nil {
		if *pu.CategoryId < 0 {
			return nil, fmt.Errorf("%w: user id must be positive", ErrInvalidInput)
		}
	}
	if pu.Price != nil {
		if pu.Price.IsNegative() {
			return nil, fmt.Errorf("%w: price must be positive", ErrInvalidInput)
		}
	}

	product, err := s.repo.Update(ctx, id, pu)
	if err != nil {
		return nil, fmt.Errorf("%w: product", translateErr(err))
	}
	return product, nil
}

func (s *productService) Delete(ctx context.Context, id int) error {

	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: user", translateErr(err))
	}
	return nil
}
