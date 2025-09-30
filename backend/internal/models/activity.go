// File: internal/models/activity.go
package models

import "time"

// Activity represents a generic timeline item.
type Activity struct {
	Type        string    `json:"type"` // "Task", "Note", "Event", "CommLog"
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	// We can add the original object if the UI needs more details
	// Data      interface{} `json:"data"`
}