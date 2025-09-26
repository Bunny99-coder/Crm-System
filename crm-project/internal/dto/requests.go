// Replace the entire contents of your internal/dto/requests.go file with this.
package dto

import "github.com/golang-jwt/jwt/v5" // <-- ADD THIS IMPORT

// --- User Request DTOs ---

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email"    validate:"required,email"`
	RoleID   int    `json:"role_id"  validate:"required,oneof=1 2"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email"    validate:"required,email"`
	RoleID   int    `json:"role_id"  validate:"required,oneof=1 2"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
	RoleID   int    `json:"role_id"  validate:"required,oneof=1 2"`
}

// --- JWT Claims DTO ---
// THIS WAS THE MISSING PIECE.
type Claims struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
	Username string `json:"username"` 

	jwt.RegisteredClaims
}