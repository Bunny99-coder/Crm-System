// File: internal/models/lead.go
package models

import "time"

// Lead represents a potential sales opportunity.
type Lead struct {
	ID          int        `db:"lead_id"       json:"id"`
	ContactID   int        `db:"contact_id"    json:"contact_id"`
	PropertyID  *int       `db:"property_id"   json:"property_id,omitempty"` // Nullable
	SourceID    int        `db:"source_id"     json:"source_id"`
	StatusID    int        `db:"status_id"     json:"status_id"`
	AssignedTo  int        `db:"assigned_to"   json:"assigned_to"`
	Notes       *string    `db:"notes"         json:"notes,omitempty"` // Nullable
	CreatedAt   time.Time  `db:"created_at"    json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"    json:"updated_at"`
}