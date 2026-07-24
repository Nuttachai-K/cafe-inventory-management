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

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetAll(ctx context.Context, limit int) ([]model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, id int, uu *model.UserUpdate) (*model.User, error)
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	err := u.db.QueryRow(
		ctx,
		`SELECT
			id,
			username,
			email,
			password_hash,
			user_role,
			created_at,
			updated_at
		FROM users
		WHERE id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.UserRole,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (u *userRepository) GetAll(ctx context.Context, limit int) ([]model.User, error) {
	query := `
	SELECT
		id,
		username,
		email,
		password_hash,
		user_role,
		created_at,
		updated_at
	FROM users 
	WHERE 1=1
	`

	args := []any{}
	argID := 1

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argID)
		args = append(args, limit)
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.UserRole,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userRepository) Create(ctx context.Context, user *model.User) error {
	return u.db.QueryRow(
		ctx,
		`INSERT INTO users(
			username,
			email,
			password_hash,
			user_role
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.UserRole,
	).Scan(&user.ID)
}

func (u *userRepository) Update(ctx context.Context, id int, uu *model.UserUpdate) (*model.User, error) {
	setClauses := []string{}
	args := []interface{}{}
	argPos := 2

	addClauses := func(column string, value interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argPos))
		args = append(args, value)
		argPos++
	}

	if uu.Username != nil {
		addClauses("username", *uu.Username)
	}
	if uu.Email != nil {
		addClauses("email", *uu.Email)
	}
	if uu.PasswordHash != nil {
		addClauses("password_hash", *uu.PasswordHash)
	}
	if uu.UserRole != nil {
		addClauses("user_role", *uu.UserRole)
	}
	if len(setClauses) == 0 {
		return nil, errors.New("no fields to update")
	}

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $1",
		strings.Join(setClauses, ", "),
	)
	args = append([]interface{}{id}, args...)

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	var user model.User
	err = u.db.QueryRow(
		ctx,
		`SELECT
			id,
			username,
			email,
			password_hash,
			user_role,
			created_at,
			updated_at
		FROM users
		WHERE id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.UserRole,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) Delete(ctx context.Context, id int) error {
	result, err := u.db.Exec(
		ctx,
		`DELETE FROM users
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
