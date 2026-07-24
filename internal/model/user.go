package model

import "time"

// User represent a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	UserRole     UserRole  `json:"user_role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserUpdate struct {
	ID           *int      `json:"id,omitempty"`
	Username     *string   `json:"username,omitempty"`
	Email        *string   `json:"email,omitempty"`
	PasswordHash *string   `json:"password_hash,omitempty"`
	UserRole     *UserRole `json:"user_role,omitempty"`
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
