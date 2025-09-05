// File: internal/models/property.go
package models

import "time"

// Property represents a real estate property in the system.
type Property struct {
	ID             int       `db:"property_id"     json:"id"`
	Name           string    `db:"name"            json:"name"`
	SiteID         int       `db:"site_id"         json:"site_id"`
	PropertyTypeID int       `db:"property_type_id" json:"property_type_id"`
	UnitNo         *string   `db:"unit_no"         json:"unit_no,omitempty"`
	SizeSqft       *float64  `db:"size_sqft"       json:"size_sqft,omitempty"` // Use float for decimal values
	Price          float64   `db:"price"           json:"price"`
	Status         string    `db:"status"          json:"status"`
	CreatedAt      time.Time `db:"created_at"      json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"      json:"updated_at"`
}