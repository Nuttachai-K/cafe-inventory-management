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

// CafeResoisitory defines database operation for cafe
type CafeRepository interface {
	GetByID(ctx context.Context, id int) (*model.Cafe, error)
	GetAll(ctx context.Context, filter model.CafeFilter) ([]model.Cafe, error)
	Create(ctx context.Context, cafe *model.Cafe) error
	Update(ctx context.Context, id int, cu *model.CafeUpdate) (*model.Cafe, error)
	Delete(ctx context.Context, id int) error
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

func (c *cafeRepository) GetAll(ctx context.Context, filter model.CafeFilter) ([]model.Cafe, error) {

	withDistance := filter.Lat != nil && filter.Lng != nil

	args := []any{}
	argID := 1
	var query string

	if withDistance {
		query = `
		SELECT id, name, address, latitude, longitude, nearest_station,
			   opening_time, closing_time, created_at, updated_at, distance_km
		FROM (
			SELECT id, name, address, latitude, longitude, nearest_station,
				   opening_time, closing_time, created_at, updated_at,
				   6371 * acos(
				 	LEAST (1.0, GREATEST(-1.0,
						cos(radians($1)) * cos(radians(latitude)) * 
						cos(radians(longitude) -  radians($2)) + 
						sin(radians($1)) * sin(radians(latitude))
					))  
				   ) AS distance_km
				FROM cafes
		) c
		WHERE 1 = 1
		`

		args = append(args, *filter.Lat, *filter.Lng)
		argID = 3
	} else {
		query = `
		SELECT id, name, address, latitude, longitude, nearest_station,
		       opening_time, closing_time, created_at, updated_at
		FROM cafes
		WHERE 1=1
		`
	}

	if filter.Station != nil {
		query += fmt.Sprintf(" AND nearest_station ILIKE $%d", argID)
		args = append(args, "%"+*filter.Station+"%")
		argID++
	}

	if withDistance && filter.RadiusKm != nil {
		query += fmt.Sprintf(" AND distance_km <= $%d", argID)
		args = append(args, *filter.RadiusKm)
		argID++
	}

	if withDistance {
		query += " ORDER BY distance_km"
	} else {
		query += " ORDER BY id"
	}

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argID)
		args = append(args, filter.Limit)
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cafes := []model.Cafe{}

	for rows.Next() {
		var cafe model.Cafe
		dest := []any{
			&cafe.ID, &cafe.Name, &cafe.Address, &cafe.Latitude, &cafe.Longitude,
			&cafe.NearestStation, &cafe.OpeningTime, &cafe.ClosingTime,
			&cafe.CreatedAt, &cafe.UpdatedAt,
		}
		var d float64
		if withDistance {
			dest = append(dest, &d)
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		if withDistance {
			cafe.DistanceKm = &d
		}
		cafes = append(cafes, cafe)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cafes, nil
}

func (c *cafeRepository) Create(ctx context.Context, cafe *model.Cafe) error {
	return c.db.QueryRow(
		ctx,
		`INSERT INTO cafes(
			name,
			address,
			latitude,
			longitude,
			nearest_station,
			opening_time,
			closing_time
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
		`,
		cafe.Name,
		cafe.Address,
		cafe.Latitude,
		cafe.Longitude,
		cafe.NearestStation,
		cafe.OpeningTime,
		cafe.ClosingTime,
	).Scan(&cafe.ID)
}

func (c *cafeRepository) Update(ctx context.Context, id int, cu *model.CafeUpdate) (*model.Cafe, error) {

	setClauses := []string{}
	args := []any{}
	argPos := 2

	addClauses := func(column string, value interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argPos))
		args = append(args, value)
		argPos++
	}

	if cu.Name != nil {
		addClauses("name", *cu.Name)
	}
	if cu.Address != nil {
		addClauses("address", *cu.Address)
	}
	if cu.Latitude != nil {
		addClauses("latitude", *cu.Latitude)
	}
	if cu.Longitude != nil {
		addClauses("longitude", *cu.Longitude)
	}
	if cu.NearestStation != nil {
		addClauses("nearest_station", *cu.NearestStation)
	}
	if cu.OpeningTime != nil {
		addClauses("opening_time", *cu.OpeningTime)
	}
	if cu.ClosingTime != nil {
		addClauses("closing_time", *cu.ClosingTime)
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no fields to update")
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(
		"UPDATE cafes SET %s WHERE id = $1",
		strings.Join(setClauses, ", "),
	)
	args = append([]interface{}{id}, args...)

	result, err := c.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	var cafe model.Cafe
	err = c.db.QueryRow(
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
		WHERE id = $1
		ORDER BY id`,
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

	return &cafe, nil
}

func (c *cafeRepository) Delete(ctx context.Context, id int) error {
	result, err := c.db.Exec(
		ctx,
		`DELETE FROM cafes
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
