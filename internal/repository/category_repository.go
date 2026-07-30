package repository

import (
	"context"
	"fmt"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	GetByID(ctx context.Context, id int) (*model.Category, error)
	GetAll(ctx context.Context, limit int) ([]model.Category, error)
	Create(ctx context.Context, cafe *model.Category) error
	Update(ctx context.Context, id int, name string) (*model.Category, error)
	Delete(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (c *categoryRepository) GetByID(ctx context.Context, id int) (*model.Category, error) {

	var category model.Category
	err := c.db.QueryRow(
		ctx,
		`SELECT
			id,
			name,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1`,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &category, err
}

func (c *categoryRepository) GetAll(
	ctx context.Context,
	limit int,
) ([]model.Category, error) {

	query := `
	SELECT
		id,
		name,
		created_at,
		updated_at
	FROM categories
	WHERE 1=1
	`

	args := []any{}
	argID := 1

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argID)
		args = append(args, limit)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category

	for rows.Next() {
		var category model.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	return c.db.QueryRow(
		ctx,
		`INSERT INTO categories(
			name
		)
		VALUES ($1)
		RETURNING id
		`,
		category.Name,
	).Scan(&category.ID)
}

func (c *categoryRepository) Update(ctx context.Context, id int, name string) (*model.Category, error) {

	result, err := c.db.Exec(
		ctx,
		"UPDATE categories SET name = $1 WHERE id = $2",
		name, id,
	)

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	var category model.Category
	err = c.db.QueryRow(
		ctx,
		`SELECT
			id,
			name,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1`,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (c *categoryRepository) Delete(ctx context.Context, id int) error {
	result, err := c.db.Exec(
		ctx,
		`DELETE FROM categories
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
