package repository

import (
	"context"

	"github.com/Nuttachai-K/cafe-finder-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CafeResoisitory defines database operation for cafe
type CafeRepository interface {
	GetByID(ctx context.Context, id int) (*model.Cafe, error)
	GetAll(ctx context.Context) ([]model.Cafe, error)
	Create(ctx context.Context, cafe *model.Cafe) error
}

type cafeRepository struct {
	db *pgxpool.Pool
}

func NewCafeRepository(db *pgxpool.Pool) CafeRepository {
	return &cafeRepository{
		db: db,
	}
}

func (c *cafeRepository) GetByID(ctx context.Context, id int) (*model.Cafe, error) {

	var cafe model.Cafe
	err := c.db.QueryRow(
		ctx,
		`SELECT
			id,
			name,
			address,
			latitude,
			longitude,
			nearest_station,
			opening_time,
			closing_time,
			created_at,
			updated_at
		FROM cafes
		WHERE id = $1`,
		id,
	).Scan(
		&cafe.ID,
		&cafe.Name,
		&cafe.Address,
		&cafe.Latitude,
		&cafe.Longitude,
		&cafe.NearestStation,
		&cafe.OpeningTime,
		&cafe.ClosingTime,
		&cafe.CreatedAt,
		&cafe.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &cafe, err
}

func (c *cafeRepository) GetAll(ctx context.Context) ([]model.Cafe, error) {
	rows, err := c.db.Query(
		ctx,
		"SELECT id, name FROM cafes",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cafes []model.Cafe

	for rows.Next() {
		var cafe model.Cafe

		err := rows.Scan(
			&cafe.ID,
			&cafe.Name,
		)
		if err != nil {
			return nil, err
		}

		cafes = append(cafes, cafe)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cafes, nil
}

func (c *cafeRepository) Create(ctx context.Context, cafe *model.Cafe) error {
	_, err := c.db.Exec(
		ctx,
		`INSERT INTO cafes(
			name,
			address,
			latitude,
			longtitude,
			nearest_station,
			opening_time,
			closing_time
		)
		VALUES ($1, $2, $3, $4)	
		`,
		cafe.Name,
		cafe.Address,
		cafe.Latitude,
		cafe.Longitude,
		cafe.NearestStation,
		cafe.OpeningTime,
		cafe.ClosingTime,
	)
	return err
}
