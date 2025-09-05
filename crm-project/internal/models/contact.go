// File: internal/models/contact.go
package models

import "time"

// Contact represents a contact in the CRM system.
type Contact struct {
	ID             int       `db:"contact_id"      json:"id"`
	FirstName      string    `db:"first_name"      json:"first_name"`
	LastName       string    `db:"last_name"       json:"last_name"`
	Email          *string   `db:"email"           json:"email,omitempty"` // omitempty means it won't appear in JSON if nil
	PrimaryPhone   string    `db:"primary_phone"   json:"primary_phone"`
	SecondaryPhone *string   `db:"secondary_phone" json:"secondary_phone,omitempty"`
	Address        *string   `db:"address"         json:"address,omitempty"`
	City           *string   `db:"city"            json:"city,omitempty"`
	SubCity        *string   `db:"sub_city"        json:"sub_city,omitempty"`
	ContactSource  *string   `db:"contact_source"  json:"contact_source,omitempty"`
	CreatedAt      time.Time `db:"created_at"      json:"created_at"`
}