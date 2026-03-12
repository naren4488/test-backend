package model

import "time"

// User is the domain model for a user (matches DB row).
// ID is a UUID. PasswordHash is never serialized to JSON.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateUserRequest is the body for POST /api/v1/users.
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Name     string `json:"name" validate:"required,max=255"`
	Password string `json:"password" validate:"required,min=6"`
}

// UpdateUserRequest is the body for PUT /api/v1/users/:id (all optional).
type UpdateUserRequest struct {
	Name  *string `json:"name" validate:"omitempty,max=255"`
	Email *string `json:"email" validate:"omitempty,email,max=255"`
}

// ResetPasswordRequest is the body for POST /api/v1/users/:id/reset-password.
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=6"`
}
