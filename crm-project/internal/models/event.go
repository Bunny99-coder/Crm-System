// File: internal/models/event.go
package models

import "time"

// Event represents a scheduled event in the calendar.
type Event struct {
	ID               int        `db:"event_id"           json:"id"`
	EventName        string     `db:"event_name"         json:"event_name"`
	EventDescription *string    `db:"event_description"  json:"event_description,omitempty"`
	StartTime        time.Time  `db:"start_time"         json:"start_time"`
	EndTime          time.Time  `db:"end_time"           json:"end_time"`
	Location         *string    `db:"location"           json:"location,omitempty"`
	OrganizerID      int        `db:"organizer_id"       json:"organizer_id"`
	CreatedAt        time.Time  `db:"created_at"         json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"         json:"updated_at"`
}