package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

type CategoryService interface {
	GetByID(ctx context.Context, id int) (*model.Category, error)
	GetAll(ctx context.Context, limit int) ([]model.Category, error)
	Create(ctx context.Context, user *model.Category) error
	Update(ctx context.Context, id int, name string) (*model.Category, error)
	Delete(ctx context.Context, id int) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) GetByID(ctx context.Context, id int) (*model.Category, error) {

	if id < 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: category", translateErr(err))
	}
	return category, nil
}

func (s *categoryService) GetAll(ctx context.Context, limit int) ([]model.Category, error) {

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.GetAll(ctx, limit)
}

func (s *categoryService) Create(ctx context.Context, category *model.Category) error {

	err := s.repo.Create(ctx, category)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: %s", ErrDuplicateCategory, category.Name)
		}
		return fmt.Errorf("%w: category", translateErr(err))
	}
	return nil
}

func (s *categoryService) Update(ctx context.Context, id int, name string) (*model.Category, error) {

	if id < 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	category, err := s.repo.Update(ctx, id, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateCategory, name)
		}
		return nil, fmt.Errorf("%w: category", translateErr(err))
	}
	return category, nil
}

func (s *categoryService) Delete(ctx context.Context, id int) error {

	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return translateErr(err)
	}
	return nil
}
