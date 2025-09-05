// Create new file: internal/repository/postgres/comm_log_repo.go
package postgres

import (
	"context"
	"crm-project/internal/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type CommLogRepo struct {
	db *sqlx.DB
}

func NewCommLogRepo(db *sqlx.DB) *CommLogRepo {
	return &CommLogRepo{db: db}
}

func (r *CommLogRepo) Create(ctx context.Context, cl models.CommLog) (int, error) {
	var newID int
	query := `INSERT INTO communication_logs (contact_id, user_id, interaction_date, interaction_type, notes)
			  VALUES ($1, $2, $3, $4, $5) RETURNING log_id`
	err := r.db.QueryRowxContext(ctx, query, cl.ContactID, cl.UserID, cl.InteractionDate, cl.InteractionType, cl.Notes).Scan(&newID)
	return newID, err
}

func (r *CommLogRepo) GetAllForContact(ctx context.Context, contactID int) ([]models.CommLog, error) {
	var logs []models.CommLog
	query := `SELECT * FROM communication_logs WHERE contact_id = $1 ORDER BY interaction_date DESC`
	err := r.db.SelectContext(ctx, &logs, query, contactID)
	return logs, err
}

func (r *CommLogRepo) GetByID(ctx context.Context, id int) (*models.CommLog, error) {
	var log models.CommLog
	query := `SELECT * FROM communication_logs WHERE log_id = $1`
	err := r.db.GetContext(ctx, &log, query, id)
	if err != nil {
		if err == sql.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &log, nil
}

// Note: Update/Delete for a log entry might not be common business requirements,
// but we include them for completeness.
func (r *CommLogRepo) Update(ctx context.Context, cl models.CommLog) error {
	query := `UPDATE communication_logs SET
				notes = $1, interaction_date = $2, interaction_type = $3
			  WHERE log_id = $4`
	result, err := r.db.ExecContext(ctx, query, cl.Notes, cl.InteractionDate, cl.InteractionType, cl.ID)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return sql.ErrNoRows }
	return nil
}

func (r *CommLogRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM communication_logs WHERE log_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return sql.ErrNoRows }
	return nil
}