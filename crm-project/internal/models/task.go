// File: internal/models/task.go
package models

import "time"

// Task represents a scheduled activity for a user.
type Task struct {
	ID              int        `db:"task_id"          json:"id"`
	TaskName        string     `db:"task_name"        json:"task_name"`
	TaskDescription *string    `db:"task_description" json:"task_description,omitempty"`
	DueDate         time.Time  `db:"due_date"         json:"due_date"`
	Status          string     `db:"status"           json:"status"`
	AssignedTo      int        `db:"assigned_to"      json:"assigned_to"`
	CreatedAt       time.Time  `db:"created_at"       json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"       json:"updated_at"`
}