// Create new file: internal/repository/postgres/lead_repo.go
package postgres

import (
	"crm-project/internal/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"context"
	"time"
)

type LeadRepo struct {
	db *sqlx.DB
}

func NewLeadRepo(db *sqlx.DB) *LeadRepo {
	return &LeadRepo{db: db}
}

func (r *LeadRepo) Create(ctx context.Context, l models.Lead) (int, error) {
	var newID int
	query := `INSERT INTO leads (contact_id, property_id, source_id, status_id, assigned_to, notes)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING lead_id`
	err := r.db.QueryRowxContext(ctx, query, l.ContactID, l.PropertyID, l.SourceID, l.StatusID, l.AssignedTo, l.Notes).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (r *LeadRepo) GetAll(ctx context.Context,) ([]models.Lead, error) {
	var leads []models.Lead
	query := `SELECT * FROM leads ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &leads, query)
	return leads, err
}

func (r *LeadRepo) GetByID(ctx context.Context, id int) (*models.Lead, error) {
	var lead models.Lead
	query := `SELECT * FROM leads WHERE lead_id = $1`
	err := r.db.GetContext(ctx, &lead, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &lead, nil
}

func (r *LeadRepo) Update(ctx context.Context, l models.Lead) error {
	query := `UPDATE leads SET
				contact_id = $1,
				property_id = $2,
				source_id = $3,
				status_id = $4,
				assigned_to = $5,
				notes = $6,
				updated_at = NOW()
			  WHERE lead_id = $7`
	result, err := r.db.ExecContext(ctx, query, l.ContactID, l.PropertyID, l.SourceID, l.StatusID, l.AssignedTo, l.Notes, l.ID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *LeadRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM leads WHERE lead_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}




// Add this struct and method to lead_repo.go
// You will need to add an import for "database/sql"

// LeadStatusCounts holds the result of our specific counting query.
type LeadStatusCounts struct {
	New       int `db:"new"`
	Contacted int `db:"contacted"`
	Qualified int `db:"qualified"`
	Converted int `db:"converted"`
	Lost      int `db:"lost"`
}

// GetLeadCountsByUserID aggregates lead counts for a specific user.
func (r *LeadRepo) GetLeadCountsByUserID(ctx context.Context, userID int) (*LeadStatusCounts, error) {
	var counts LeadStatusCounts
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN ls.name = 'New' THEN 1 ELSE 0 END), 0) AS new,
			COALESCE(SUM(CASE WHEN ls.name = 'Contacted' THEN 1 ELSE 0 END), 0) AS contacted,
			COALESCE(SUM(CASE WHEN ls.name = 'Qualified' THEN 1 ELSE 0 END), 0) AS qualified,
			COALESCE(SUM(CASE WHEN ls.name = 'Converted' THEN 1 ELSE 0 END), 0) AS converted,
			COALESCE(SUM(CASE WHEN ls.name = 'Lost' THEN 1 ELSE 0 END), 0) AS lost
		FROM leads l
		JOIN lead_statuses ls ON l.status_id = ls.status_id
		WHERE l.assigned_to = $1
	`
	// Use the context-aware GetContext method
	err := r.db.GetContext(ctx, &counts, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &LeadStatusCounts{}, nil
		}
		return nil, err
	}
	return &counts, nil
}





// Add this new struct to lead_repo.go to hold the rich report data
type SourceLeadReportRow struct {
	LeadDate         time.Time `db:"lead_date" json:"lead_date"`
	ContactName      string    `db:"contact_name" json:"contact_name"`
	ContactPhone     string    `db:"contact_phone" json:"contact_phone"`
	ContactEmail     *string   `db:"contact_email" json:"contact_email"`
	LeadSource       string    `db:"lead_source" json:"lead_source"`
	AssignedEmployee string    `db:"assigned_employee" json:"assigned_employee"`
	LeadStatus       string    `db:"lead_status" json:"lead_status"`
	// updated_by is not in our schema, so we'll omit it for now.
}

// GetSourceLeadReport retrieves a detailed list of leads, joining with other tables.
func (r *LeadRepo) GetSourceLeadReport(ctx context.Context) ([]SourceLeadReportRow, error) {
	var reportRows []SourceLeadReportRow
	query := `
		SELECT
			l.created_at AS lead_date,
			c.first_name || ' ' || c.last_name AS contact_name,
			c.primary_phone AS contact_phone,
			c.email AS contact_email,
			ls.name AS lead_source,
			u.username AS assigned_employee,
			lst.name AS lead_status
		FROM
			leads l
		JOIN
			contacts c ON l.contact_id = c.contact_id
		JOIN
			lead_sources ls ON l.source_id = ls.source_id
		JOIN
			users u ON l.assigned_to = u.user_id
		JOIN
			lead_statuses lst ON l.status_id = lst.status_id
		ORDER BY
			l.created_at DESC
	`
	// In a real app, you would add WHERE clauses here for date filtering.
	err := r.db.SelectContext(ctx, &reportRows, query)
	return reportRows, err
}



// in lead_repo.go
func (r *LeadRepo) CheckForOpenLeadByContactID(ctx context.Context, contactID int) (bool, error) {
    var exists bool
    query := `
        SELECT EXISTS (
            SELECT 1 FROM leads l
            JOIN lead_statuses ls ON l.status_id = ls.status_id
            WHERE l.contact_id = $1 AND ls.name NOT IN ('Converted', 'Lost')
        )
    `
    err := r.db.GetContext(ctx, &exists, query, contactID)
    return exists, err
}

func (r *LeadRepo) GetAllLeadsForUser(ctx context.Context, userID int) ([]models.Lead, error) {
	var leads []models.Lead
	query := `SELECT * FROM leads WHERE assigned_to = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &leads, query, userID)
	return leads, err
}