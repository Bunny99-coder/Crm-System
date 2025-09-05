// Create new file: internal/repository/postgres/event_repo.go
package postgres

import (
	"context"
	"crm-project/internal/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type EventRepo struct {
	db *sqlx.DB
}

func NewEventRepo(db *sqlx.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Create(ctx context.Context, e models.Event) (int, error) {
	var newID int
	query := `INSERT INTO events (event_name, event_description, start_time, end_time, location, organizer_id)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING event_id`
	err := r.db.QueryRowxContext(ctx, query, e.EventName, e.EventDescription, e.StartTime, e.EndTime, e.Location, e.OrganizerID).Scan(&newID)
	return newID, err
}

// in event_repo.go
func (r *EventRepo) GetAll(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	query := `SELECT * FROM events ORDER BY start_time ASC`
	err := r.db.SelectContext(ctx, &events, query)
	return events, err
}


// GetAllForUser retrieves all events organized by a specific user.
func (r *EventRepo) GetAllForUser(ctx context.Context, organizerID int) ([]models.Event, error) {
	var events []models.Event
	query := `SELECT * FROM events WHERE organizer_id = $1 ORDER BY start_time ASC`
	err := r.db.SelectContext(ctx, &events, query, organizerID)
	return events, err
}

func (r *EventRepo) GetByID(ctx context.Context, id int) (*models.Event, error) {
	var event models.Event
	query := `SELECT * FROM events WHERE event_id = $1`
	err := r.db.GetContext(ctx, &event, query, id)
	if err != nil {
		if err == sql.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &event, nil
}

func (r *EventRepo) Update(ctx context.Context, e models.Event) error {
	query := `UPDATE events SET
				event_name = $1, event_description = $2, start_time = $3, end_time = $4,
				location = $5, organizer_id = $6, updated_at = NOW()
			  WHERE event_id = $7`
	result, err := r.db.ExecContext(ctx, query, e.EventName, e.EventDescription, e.StartTime, e.EndTime, e.Location, e.OrganizerID, e.ID)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return sql.ErrNoRows }
	return nil
}

func (r *EventRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM events WHERE event_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil { return err }
	rows, _ := result.RowsAffected()
	if rows == 0 { return sql.ErrNoRows }
	return nil
}