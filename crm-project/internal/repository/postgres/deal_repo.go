// Create new file: internal/repository/postgres/deal_repo.go
package postgres

import (
	"crm-project/internal/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"context"
)

type DealRepo struct {
	db *sqlx.DB
}

func NewDealRepo(db *sqlx.DB) *DealRepo {
	return &DealRepo{db: db}
}

func (r *DealRepo) Create(ctx context.Context,d models.Deal) (int, error) {
	var newID int
	query := `INSERT INTO deals (lead_id, property_id, stage_id, deal_status, deal_amount, closing_date, notes)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING deal_id`
	err := r.db.QueryRowxContext(ctx,query, d.LeadID, d.PropertyID, d.StageID, d.DealStatus, d.DealAmount, d.ClosingDate, d.Notes).Scan(&newID)
	return newID, err
}

func (r *DealRepo) GetAll(ctx context.Context,) ([]models.Deal, error) {
	var deals []models.Deal
	query := `SELECT * FROM deals ORDER BY deal_date DESC`
	err := r.db.SelectContext(ctx,&deals, query)
	return deals, err
}

func (r *DealRepo) GetByID( ctx context.Context,id int) (*models.Deal, error) {
	var deal models.Deal
	query := `SELECT * FROM deals WHERE deal_id = $1`
	err := r.db.GetContext(ctx, &deal, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &deal, nil
}

func (r *DealRepo) Update(ctx context.Context, d models.Deal) error {
	query := `UPDATE deals SET
				lead_id = $1,
				property_id = $2,
				stage_id = $3,
				deal_status = $4,
				deal_amount = $5,
				closing_date = $6,
				notes = $7,
				updated_at = NOW()
			  WHERE deal_id = $8`
	result, err := r.db.ExecContext(ctx, query, d.LeadID, d.PropertyID, d.StageID, d.DealStatus, d.DealAmount, d.ClosingDate, d.Notes, d.ID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *DealRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM deals WHERE deal_id = $1`
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


// Add this struct to deal_repo.go
type EmployeeSalesReportRow struct {
	EmployeeName    string  `db:"employee_name" json:"employee_name"`
	NumberOfSales   int     `db:"number_of_sales" json:"number_of_sales"`
	TotalSalesAmount float64 `db:"total_sales_amount" json:"total_sales_amount"`
}

// GetEmployeeSalesReport aggregates sales data per employee.
func (r *DealRepo) GetEmployeeSalesReport(ctx context.Context) ([]EmployeeSalesReportRow, error) {
	var reportRows []EmployeeSalesReportRow
	query := `
		SELECT
			u.username AS employee_name,
			COUNT(d.deal_id) AS number_of_sales,
			COALESCE(SUM(d.deal_amount), 0) AS total_sales_amount
		FROM
			deals d
		JOIN
			leads l ON d.lead_id = l.lead_id
		JOIN
			users u ON l.assigned_to = u.user_id
		WHERE
			d.deal_status = 'Closed-Won' -- Only count successful deals
		GROUP BY
			u.username
		ORDER BY
			total_sales_amount DESC
	`
	err := r.db.SelectContext(ctx, &reportRows, query)
	return reportRows, err
}





// Add this struct to deal_repo.go
type SourceSalesReportRow struct {
	SourceName       string  `db:"source_name" json:"source_name"`
	NumberOfSales    int     `db:"number_of_sales" json:"number_of_sales"`
	TotalSalesAmount float64 `db:"total_sales_amount" json:"total_sales_amount"`
}

// GetSourceSalesReport aggregates sales data per lead source.
func (r *DealRepo) GetSourceSalesReport(ctx context.Context) ([]SourceSalesReportRow, error) {
	var reportRows []SourceSalesReportRow
	query := `
		SELECT
			ls.name AS source_name,
			COUNT(d.deal_id) AS number_of_sales,
			COALESCE(SUM(d.deal_amount), 0) AS total_sales_amount
		FROM
			deals d
		JOIN
			leads l ON d.lead_id = l.lead_id
		JOIN
			lead_sources ls ON l.source_id = ls.source_id
		WHERE
			d.deal_status = 'Closed-Won'
		GROUP BY
			ls.name
		ORDER BY
			total_sales_amount DESC
	`
	err := r.db.SelectContext(ctx, &reportRows, query)
	return reportRows, err
}


// in lead_repo.go
func (r *LeadRepo) GetAllForUser(ctx context.Context, userID int) ([]models.Lead, error) {
	var leads []models.Lead
	query := `SELECT * FROM leads WHERE assigned_to = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &leads, query, userID)
	return leads, err
}