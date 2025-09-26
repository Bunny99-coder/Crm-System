package models

import "time"

// Task represents a scheduled activity for a user, potentially linked to a lead or deal.
type Task struct {
    ID              int        `db:"task_id"          json:"id"`
    TaskName        string     `db:"task_name"        json:"task_name"`
    TaskDescription *string    `db:"task_description" json:"task_description,omitempty"`
    DueDate         time.Time  `db:"due_date"         json:"due_date"`
    Status          string     `db:"status"           json:"status"`
    AssignedTo      int        `db:"assigned_to"      json:"assigned_to"`
    LeadID          *int       `db:"lead_id"          json:"lead_id,omitempty"`
    DealID          *int       `db:"deal_id"          json:"deal_id,omitempty"`
    CreatedAt       time.Time  `db:"created_at"       json:"created_at"`
    UpdatedAt       *time.Time `db:"updated_at"       json:"updated_at,omitempty"`
    DeletedAt       *time.Time `db:"deleted_at"       json:"deleted_at,omitempty"`
    CreatedBy       int        `db:"created_by"       json:"created_by"`
}