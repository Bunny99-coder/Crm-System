// File: internal/models/user.go
package models

import "time"

// User represents a user in the system.
type User struct {
	ID           int       `db:"user_id"      json:"id"`
	Username     string    `db:"username"     json:"username"`
	Password     string    `db:"-"            json:"password,omitempty"` // Used for input, but not stored in DB
	PasswordHash string    `db:"password_hash" json:"-"`                 // Stored in DB, but never sent in JSON responses
	Email        string    `db:"email"        json:"email"`
	RoleID       int       `db:"role_id"      json:"role_id"`
	CreatedAt    time.Time `db:"created_at"   json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"   json:"updated_at"`
	RoleName     string    `db:"-"            json:"role_name"` // Populated from roles table
}