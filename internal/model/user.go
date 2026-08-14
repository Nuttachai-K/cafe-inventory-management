package model

import (
	"encoding/json"
	"time"
)

// User represent a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"-"`
	UserRole     UserRole  `json:"user_role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserUpdate struct {
	ID           *int      `json:"id,omitempty"`
	Username     *string   `json:"username,omitempty"`
	Email        *string   `json:"email,omitempty"`
	Password     *string   `json:"password,omitempty"`
	PasswordHash *string   `json:"password_hash,omitempty"`
	UserRole     *UserRole `json:"user_role,omitempty"`
	IsActive     *bool     `json:"is_active,omitempty"`
}

// For swagger documentation
type UserCreate struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	UserRole UserRole `json:"user_role"`
}

type UserRole string

const (
	RoleAdmin UserRole = "ADMIN"
	RoleStaff UserRole = "STAFF"
)

func (r UserRole) Valid() bool {
	switch r {
	case RoleAdmin, RoleStaff:
		return true
	default:
		return false
	}
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(&u),
		CreatedAt: formatDateHourJST(u.CreatedAt),
		UpdatedAt: formatDateHourJST(u.UpdatedAt),
	})
}
