package models

import (
	"database/sql"
	"time"
)

// Deal represents a finalized sales deal in the system.
type Deal struct {
	ID          int           `db:"deal_id"      json:"id"`
	LeadID      int           `db:"lead_id"      json:"lead_id"`
	PropertyID  int           `db:"property_id"  json:"property_id"`
	StageID     int           `db:"stage_id"     json:"stage_id"`
	DealStatus  string        `db:"deal_status"  json:"deal_status"`
	DealAmount  float64       `db:"deal_amount"  json:"deal_amount"`
	DealDate    time.Time     `db:"deal_date"    json:"deal_date"`
	ClosingDate *time.Time    `db:"closing_date" json:"closing_date,omitempty"`
	Notes       *string       `db:"notes"        json:"notes,omitempty"`
	CreatedAt   time.Time     `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"   json:"updated_at"`
	CreatedBy   sql.NullInt64 `db:"created_by"   json:"created_by"`
}