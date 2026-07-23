package model

// User represent a user in the system
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	UserRole     string `json:"user_role"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UserUpdate struct {
	ID           *int    `json:"id,omitempty"`
	Username     *string `json:"username,omitempty"`
	Email        *string `json:"email,omitempty"`
	PasswordHash *string `json:"password_hash,omitempty"`
	UserRole     *string `json:"user_role,omitempty"`
}
