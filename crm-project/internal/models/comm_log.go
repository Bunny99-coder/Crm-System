// File: internal/models/comm_log.go
package models

import "time"

// CommLog represents an entry in the communication log.
type CommLog struct {
	ID              int       `db:"log_id"           json:"id"`
	ContactID       int       `db:"contact_id"       json:"contact_id"`
	UserID          int       `db:"user_id"          json:"user_id"`
	InteractionDate time.Time `db:"interaction_date" json:"interaction_date"`
	InteractionType string    `db:"interaction_type" json:"interaction_type"`
	Notes           *string   `db:"notes"            json:"notes,omitempty"`
	CreatedAt       time.Time `db:"created_at"       json:"created_at"`
}