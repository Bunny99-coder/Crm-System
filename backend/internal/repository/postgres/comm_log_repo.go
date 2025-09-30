package postgres

import (
    "database/sql"
    "fmt"
    "time"

    "crm-project/internal/models"
    "github.com/jmoiron/sqlx"
)

// CommLogRepo implements the CommLogRepository interface
type CommLogRepo struct {
    db *sqlx.DB
}

func NewCommLogRepository(db *sqlx.DB) *CommLogRepo {
    return &CommLogRepo{db: db}
}

// CreateCommLog creates a new communication log
func (r *CommLogRepo) CreateCommLog(log *models.CommLog) error {
    query := `
        INSERT INTO communication_logs (contact_id, user_id, lead_id, deal_id, interaction_date, interaction_type, notes, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING log_id, created_at, updated_at
    `

    currentTime := time.Now()
    var createdAt time.Time
    var updatedAt time.Time
    err := r.db.QueryRowx(
        query,
        log.ContactID,
        log.UserID,
        log.LeadID,
        log.DealID,
        log.InteractionDate,
        log.InteractionType,
        log.Notes,
        currentTime,
        currentTime,
    ).Scan(&log.ID, &createdAt, &updatedAt)

    if err != nil {
        return fmt.Errorf("failed to create communication log: %w", err)
    }

    log.CreatedAt = &createdAt
    log.UpdatedAt = &updatedAt
    return nil
}

// GetCommLogByID retrieves a communication log by ID
func (r *CommLogRepo) GetCommLogByID(id int) (*models.CommLog, error) {
    query := `
        SELECT log_id, contact_id, user_id, lead_id, deal_id, interaction_date, interaction_type, notes, created_at, updated_at, deleted_at
        FROM communication_logs
        WHERE log_id = $1 AND deleted_at IS NULL
    `

    var log models.CommLog
    err := r.db.Get(&log, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("communication log not found")
        }
        return nil, fmt.Errorf("failed to get communication log: %w", err)
    }

    return &log, nil
}

// GetCommLogsByDealID retrieves all communication logs for a specific deal
func (r *CommLogRepo) GetCommLogsByDealID(dealID int) ([]models.CommLog, error) {
    query := `
        SELECT log_id, contact_id, user_id, lead_id, deal_id, interaction_date, interaction_type, notes, created_at, deleted_at
        FROM communication_logs
        WHERE deal_id = $1 AND deleted_at IS NULL
        ORDER BY interaction_date DESC
    `

    var logs []models.CommLog
    err := r.db.Select(&logs, query, dealID)
    if err != nil {
        return nil, fmt.Errorf("failed to get communication logs by deal ID: %w", err)
    }

    return logs, nil
}

// GetCommLogsByContactID retrieves all communication logs for a specific contact
func (r *CommLogRepo) GetCommLogsByContactID(contactID int) ([]models.CommLog, error) {
    query := `
        SELECT log_id, contact_id, user_id, lead_id, deal_id, interaction_date, interaction_type, notes, created_at, deleted_at
        FROM communication_logs
        WHERE contact_id = $1 AND deleted_at IS NULL
        ORDER BY interaction_date DESC
    `

    var logs []models.CommLog
    err := r.db.Select(&logs, query, contactID)
    if err != nil {
        return nil, fmt.Errorf("failed to get communication logs by contact ID: %w", err)
    }

    return logs, nil
}

// GetAllCommLogs retrieves all communication logs
func (r *CommLogRepo) GetAllCommLogs() ([]models.CommLog, error) {
    query := `
        SELECT log_id, contact_id, user_id, lead_id, deal_id, interaction_date, interaction_type, notes, created_at, deleted_at
        FROM communication_logs
        WHERE deleted_at IS NULL
        ORDER BY interaction_date DESC
    `

    var logs []models.CommLog
    err := r.db.Select(&logs, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get all communication logs: %w", err)
    }

    return logs, nil
}

// UpdateCommLog updates an existing communication log
func (r *CommLogRepo) UpdateCommLog(log *models.CommLog) error {
    query := `
        UPDATE communication_logs 
        SET contact_id = $1, user_id = $2, lead_id = $3, deal_id = $4, interaction_date = $5, interaction_type = $6, notes = $7, created_at = $8
        WHERE log_id = $9 AND deleted_at IS NULL
        RETURNING created_at
    `

    currentTime := time.Now()
    var createdAt time.Time
    err := r.db.QueryRowx(
        query,
        log.ContactID,
        log.UserID,
        log.LeadID,
        log.DealID,
        log.InteractionDate,
        log.InteractionType,
        log.Notes,
        currentTime,
        log.ID,
    ).Scan(&createdAt)

    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("communication log not found")
        }
        return fmt.Errorf("failed to update communication log: %w", err)
    }

    log.CreatedAt = &createdAt
    return nil
}

// DeleteCommLog soft deletes a communication log
func (r *CommLogRepo) DeleteCommLog(id int) error {
    query := `
        UPDATE communication_logs 
        SET deleted_at = $1
        WHERE log_id = $2 AND deleted_at IS NULL
    `

    result, err := r.db.Exec(query, time.Now(), id)
    if err != nil {
        return fmt.Errorf("failed to delete communication log: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("communication log not found")
    }

    return nil
}

// GetCommLogsForUser retrieves communication logs for a specific user
func (r *CommLogRepo) GetCommLogsForUser(userID int) ([]models.CommLog, error) {
    query := `
        SELECT log_id, contact_id, user_id, lead_id, deal_id, interaction_date, interaction_type, notes, created_at, deleted_at
        FROM communication_logs
        WHERE user_id = $1 AND deleted_at IS NULL
        ORDER BY interaction_date DESC
    `

    var logs []models.CommLog
    err := r.db.Select(&logs, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get communication logs for user: %w", err)
    }

    return logs, nil
}