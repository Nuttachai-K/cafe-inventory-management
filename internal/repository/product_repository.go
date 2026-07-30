package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id int) (*model.ProductWithCategory, error)
	GetAll(ctx context.Context, limit int) ([]model.ProductWithCategory, error)
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, id int, pu *model.ProductUpdate) (*model.ProductWithCategory, error)
	Delete(ctx context.Context, id int) error
}

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (p *productRepository) GetByID(ctx context.Context, id int) (*model.ProductWithCategory, error) {

	var product model.ProductWithCategory
	err := p.db.QueryRow(
		ctx,
		`SELECT
			id,
			cafe_id,
			category_id,
			name,
			description,
			price,
			created_at,
			updated_at
		FROM products
		WHERE id = $1`,
		id,
	).Scan(
		&product.ID,
		&product.CafeId,
		&product.CategoryId,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, err
}

func (p *productRepository) GetAll(ctx context.Context, limit int) ([]model.ProductWithCategory, error) {

	query := `
		SELECT
			products.id,
			products.cafe_id,
			products.category_id,
			categories.name,
			products.name,
			products.description,
			products.price,
			products.created_at,
			products.updated_at
		FROM products
		JOIN categories 
		ON products.category_id = categories.id
		WHERE 1=1
		`

	args := []any{}
	argID := 1

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argID)
		args = append(args, limit)
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.ProductWithCategory

	for rows.Next() {
		var product model.ProductWithCategory

		err := rows.Scan(
			&product.ID,
			&product.CafeId,
			&product.CategoryId,
			&product.CategoryName,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *productRepository) Create(ctx context.Context, product *model.Product) error {

	return p.db.QueryRow(
		ctx,
		`INSERT INTO products(
			cafe_id,
			category_id,
			name,
			description,
			price
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`,
		product.CafeId,
		product.CategoryId,
		product.Name,
		product.Description,
		product.Price,
	).Scan(&product.ID)
}

func (p *productRepository) Update(ctx context.Context, id int, pu *model.ProductUpdate) (*model.ProductWithCategory, error) {

	setClauses := []string{}
	args := []interface{}{}
	argPos := 2

	addClauses := func(column string, value any) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argPos))
		args = append(args, value)
		argPos++
	}
	if pu.CafeId != nil {
		addClauses("cafe_id", *pu.CafeId)
	}
	if pu.CategoryId != nil {
		addClauses("category_id", *pu.CategoryId)
	}
	if pu.Name != nil {
		addClauses("name", *pu.Name)
	}
	if pu.Description != nil {
		addClauses("description", *pu.Description)
	}
	if pu.Price != nil {
		addClauses("price", *pu.Price)
	}
	if len(setClauses) == 0 {
		return nil, errors.New("no fields to update")
	}

	query := fmt.Sprintf(
		"UPDATE products SET %s WHERE products.id = $1",
		strings.Join(setClauses, ", "),
	)
	args = append([]any{id}, args...)

	result, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	var product model.ProductWithCategory
	err = p.db.QueryRow(
		ctx,
		`SELECT
			products.id,
			products.cafe_id,
			products.category_id,
			categories.name,
			products.name,
			products.description,
			products.price,
			products.created_at,
			products.updated_at
		FROM products
		JOIN categories 
		ON products.category_id = categories.id
		WHERE products.id = $1`,
		id,
	).Scan(
		&product.ID,
		&product.CafeId,
		&product.CategoryId,
		&product.CategoryName,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *productRepository) Delete(ctx context.Context, id int) error {

	result, err := p.db.Exec(
		ctx,
		`DELETE FROM products
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}
