package postgres

import (
    "database/sql" // Add this import
    "fmt"
    "time"

    "crm-project/internal/models"
    "github.com/jmoiron/sqlx"
)

// EventRepo implements the EventRepository interface
type EventRepo struct {
    db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepo {
    return &EventRepo{db: db}
}

// CreateEvent creates a new event
func (r *EventRepo) CreateEvent(event *models.Event) error {
    query := `
        INSERT INTO events (event_name, event_description, start_time, end_time, location, organizer_id, lead_id, deal_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING event_id, updated_at
    `

    currentTime := time.Now()
    var updatedAt time.Time
    err := r.db.QueryRowx(
        query,
        event.EventName,
        event.EventDescription,
        event.StartTime,
        event.EndTime,
        event.Location,
        event.OrganizerID,
        event.LeadID,
        event.DealID,
        event.CreatedAt,
        currentTime,
    ).Scan(&event.ID, &updatedAt)

    if err != nil {
        return fmt.Errorf("failed to create event: %w", err)
    }

    event.UpdatedAt = &updatedAt
    return nil
}

// GetEventByID retrieves an event by ID
func (r *EventRepo) GetEventByID(id int) (*models.Event, error) {
    query := `
        SELECT event_id, event_name, event_description, start_time, end_time, location, organizer_id, lead_id, deal_id, created_at, updated_at, deleted_at
        FROM events
        WHERE event_id = $1 AND deleted_at IS NULL
    `

    var event models.Event
    var updatedAt sql.NullTime
    var deletedAt sql.NullTime
    err := r.db.Get(&event, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("event not found")
        }
        return nil, fmt.Errorf("failed to get event: %w", err)
    }

    if updatedAt.Valid {
        event.UpdatedAt = &updatedAt.Time
    }
    if deletedAt.Valid {
        event.DeletedAt = &deletedAt.Time
    }

    return &event, nil
}

// GetEventsByDealID retrieves all events for a specific deal
func (r *EventRepo) GetEventsByDealID(dealID int) ([]models.Event, error) {
    query := `
        SELECT event_id, event_name, event_description, start_time, end_time, location, organizer_id, lead_id, deal_id, created_at, updated_at, deleted_at
        FROM events
        WHERE deal_id = $1 AND deleted_at IS NULL
        ORDER BY start_time ASC
    `

    var events []models.Event
    err := r.db.Select(&events, query, dealID)
    if err != nil {
        return nil, fmt.Errorf("failed to get events by deal ID: %w", err)
    }

    return events, nil
}

// GetAllEvents retrieves all events (global)
func (r *EventRepo) GetAllEvents() ([]models.Event, error) {
    query := `
        SELECT event_id, event_name, event_description, start_time, end_time, location, organizer_id, lead_id, deal_id, created_at, updated_at, deleted_at
        FROM events
        WHERE deleted_at IS NULL
        ORDER BY start_time ASC
    `

    var events []models.Event
    err := r.db.Select(&events, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get all events: %w", err)
    }

    return events, nil
}

// UpdateEvent updates an existing event
func (r *EventRepo) UpdateEvent(event *models.Event) error {
    query := `
        UPDATE events 
        SET event_name = $1, event_description = $2, start_time = $3, end_time = $4, location = $5, organizer_id = $6, lead_id = $7, deal_id = $8, updated_at = $9
        WHERE event_id = $10 AND deleted_at IS NULL
        RETURNING updated_at
    `

    currentTime := time.Now()
    var updatedAt time.Time
    err := r.db.QueryRowx(
        query,
        event.EventName,
        event.EventDescription,
        event.StartTime,
        event.EndTime,
        event.Location,
        event.OrganizerID,
        event.LeadID,
        event.DealID,
        currentTime,
        event.ID,
    ).Scan(&updatedAt)

    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("event not found")
        }
        return fmt.Errorf("failed to update event: %w", err)
    }

    event.UpdatedAt = &updatedAt
    return nil
}

// DeleteEvent soft deletes an event
func (r *EventRepo) DeleteEvent(id int) error {
    query := `
        UPDATE events 
        SET deleted_at = $1
        WHERE event_id = $2 AND deleted_at IS NULL
    `

    result, err := r.db.Exec(query, time.Now(), id)
    if err != nil {
        return fmt.Errorf("failed to delete event: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("event not found")
    }

    return nil
}

// GetEventsForUser retrieves events for a specific user (organizer)
func (r *EventRepo) GetEventsForUser(userID int) ([]models.Event, error) {
    query := `
        SELECT event_id, event_name, event_description, start_time, end_time, location, organizer_id, lead_id, deal_id, created_at, updated_at, deleted_at
        FROM events
        WHERE organizer_id = $1 AND deleted_at IS NULL
        ORDER BY start_time ASC
    `

    var events []models.Event
    err := r.db.Select(&events, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get events for user: %w", err)
    }

    return events, nil
}