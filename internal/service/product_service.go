package service

import (
	"context"
	"fmt"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
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
	pool *pgxpool.Pool
}

func NewProductService(repo repository.ProductRepository, pool *pgxpool.Pool) ProductService {
	return &productService{
		repo: repo,
		pool: pool,
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

	products, err := s.repo.GetAll(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: products", translateErr(err))
	}
	return products, nil
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

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	productRepoTx := repository.NewProductRepository(tx)
	inventoryRepoTx := repository.NewInventoryRepository(tx)

	if err := productRepoTx.Create(ctx, product); err != nil {
		return err
	}

	inventory := &model.Inventory{ProductId: product.ID, StockQuantity: 0}
	if err := inventoryRepoTx.Create(ctx, inventory); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *productService) Update(ctx context.Context, id int, pu *model.ProductUpdate) (*model.ProductWithCategory, error) {

	if id < 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	if pu.CafeId != nil {
		if *pu.CafeId < 0 {
			return nil, fmt.Errorf("%w: cafe id must be positive", ErrInvalidInput)
		}
	}
	if pu.CategoryId != nil {
		if *pu.CategoryId < 0 {
			return nil, fmt.Errorf("%w: category id must be positive", ErrInvalidInput)
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

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	productRepoTx := repository.NewProductRepository(tx)
	inventoryRepoTx := repository.NewInventoryRepository(tx)

	if err := inventoryRepoTx.Delete(ctx, id); err != nil {
		return translateErr(err)
	}
	if err := productRepoTx.Delete(ctx, id); err != nil {
		return translateErr(err)
	}

	return tx.Commit(ctx)
}
