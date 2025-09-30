// internal/models/note.go
package models

import "time"

// Note represents a text note created by a user and optionally linked to a contact, lead, or deal.
// It aligns with the frontend's Note interface, but uses 'content' for the note text.
type Note struct {
	ID        int        `db:"note_id" json:"id"`         // Primary key
	UserID    int        `db:"user_id" json:"created_by"` // Who created the note
	ContactID *int       `db:"contact_id" json:"contact_id,omitempty"` // Optional link to a contact
	LeadID    *int       `db:"lead_id" json:"lead_id,omitempty"`       // Optional link to a lead
	DealID    *int       `db:"deal_id" json:"deal_id,omitempty"`       // Optional link to a deal
	Content   string     `db:"note_text" json:"content"`   // The actual note text (renamed for frontend consistency)
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"` // Use pointer for nullable updated_at
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"` // New field
}