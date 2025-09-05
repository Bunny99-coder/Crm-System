// Create new file: internal/repository/postgres/property_repo.go
package postgres

import (
	"crm-project/internal/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"context"
)

type PropertyRepo struct {
	db *sqlx.DB
}

func NewPropertyRepo(db *sqlx.DB) *PropertyRepo {
	return &PropertyRepo{db: db}
}

// Create inserts a new property into the database.
func (r *PropertyRepo) Create(ctx context.Context,p models.Property) (int, error) {
	var newID int
	query := `INSERT INTO properties (name, site_id, property_type_id, unit_no, size_sqft, price, status)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING property_id`
	err := r.db.QueryRowxContext(ctx, query, p.Name, p.SiteID, p.PropertyTypeID, p.UnitNo, p.SizeSqft, p.Price, p.Status).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

// GetAll retrieves all properties from the database.
func (r *PropertyRepo) GetAll(ctx context.Context,) ([]models.Property, error) {
	var properties []models.Property
	query := `SELECT * FROM properties ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &properties, query)
	if err != nil {
		return nil, err
	}
	return properties, nil
}

// GetByID retrieves a single property by its ID.
func (r *PropertyRepo) GetByID(ctx context.Context, id int) (*models.Property, error) {
	var property models.Property
	query := `SELECT * FROM properties WHERE property_id = $1`
	err := r.db.GetContext(ctx, &property, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &property, nil
}

// Update modifies an existing property in the database.
func (r *PropertyRepo) Update(ctx context.Context, p models.Property) error {
	query := `UPDATE properties SET
				name = $1,
				site_id = $2,
				property_type_id = $3,
				unit_no = $4,
				size_sqft = $5,
				price = $6,
				status = $7,
				updated_at = NOW()
			  WHERE property_id = $8`
	result, err := r.db.ExecContext(ctx, query, p.Name, p.SiteID, p.PropertyTypeID, p.UnitNo, p.SizeSqft, p.Price, p.Status, p.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delete removes a property from the database by its ID.
func (r *PropertyRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM properties WHERE property_id = $1`
	result, err := r.db.ExecContext(ctx,query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// in property_repo.go

// IsPropertyInOpenLeadOrDeal checks if a property is already part of an active sales process.
func (r *PropertyRepo) IsPropertyInOpenLeadOrDeal(ctx context.Context, propertyID int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM leads l
			JOIN lead_statuses ls ON l.status_id = ls.status_id
			WHERE l.property_id = $1 AND ls.name NOT IN ('Converted', 'Lost')
			UNION
			SELECT 1 FROM deals d
			WHERE d.property_id = $1 AND d.deal_status NOT IN ('Closed-Won', 'Closed-Lost')
		)
	`
	err := r.db.GetContext(ctx, &exists, query, propertyID)
	return exists, err
}