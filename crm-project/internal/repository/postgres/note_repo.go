package postgres

import (
    "database/sql"
    "fmt"
    "time"

    "crm-project/internal/models"
    "github.com/jmoiron/sqlx"
)

// NoteRepo implements NoteRepository interface
type NoteRepo struct {
    db *sqlx.DB
}

func NewNoteRepository(db *sqlx.DB) *NoteRepo {
    return &NoteRepo{db: db}
}

// CreateNote creates a new note
func (r *NoteRepo) CreateNote(note *models.Note) error {
    query := `
        INSERT INTO notes (user_id, contact_id, lead_id, deal_id, note_text, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING note_id
    `

    currentTime := time.Now()
    err := r.db.QueryRowx(
        query,
        note.UserID,
        note.ContactID,
        note.LeadID,
        note.DealID,
        note.Content,
        currentTime,
        currentTime,
    ).Scan(&note.ID)

    if err != nil {
        return fmt.Errorf("failed to create note: %w", err)
    }

    note.CreatedAt = currentTime
    note.UpdatedAt = &currentTime

    return nil
}

// GetNoteByID retrieves a note by its ID
func (r *NoteRepo) GetNoteByID(id int) (*models.Note, error) {
    query := `
        SELECT note_id, user_id, contact_id, lead_id, deal_id, note_text, created_at, updated_at, deleted_at
        FROM notes
        WHERE note_id = $1 AND deleted_at IS NULL
    `

    var note models.Note
    err := r.db.Get(&note, query, id)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("note not found")
        }
        return nil, fmt.Errorf("failed to get note: %w", err)
    }

    return &note, nil
}

// GetNotesByContactID retrieves all notes for a specific contact
func (r *NoteRepo) GetNotesByContactID(contactID int) ([]models.Note, error) {
    query := `
        SELECT note_id, user_id, contact_id, lead_id, deal_id, note_text, created_at, updated_at
        FROM notes
        WHERE contact_id = $1 AND deleted_at IS NULL
        ORDER BY created_at DESC
    `

    var notes []models.Note
    err := r.db.Select(&notes, query, contactID)

    if err != nil {
        return nil, fmt.Errorf("failed to get notes by contact ID: %w", err)
    }

    return notes, nil
}

// UpdateNote updates an existing note
func (r *NoteRepo) UpdateNote(note *models.Note) error {
    query := `
        UPDATE notes 
        SET note_text = $1, updated_at = $2
        WHERE note_id = $3 AND deleted_at IS NULL
        RETURNING updated_at
    `

    currentTime := time.Now()
    var updatedAt time.Time

    err := r.db.QueryRowx(
        query,
        note.Content,
        currentTime,
        note.ID,
    ).Scan(&updatedAt)

    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("note not found")
        }
        return fmt.Errorf("failed to update note: %w", err)
    }

    note.UpdatedAt = &updatedAt
    return nil
}

// DeleteNote soft deletes a note
func (r *NoteRepo) DeleteNote(id int) error {
    query := `
        UPDATE notes 
        SET deleted_at = $1
        WHERE note_id = $2 AND deleted_at IS NULL
    `

    result, err := r.db.Exec(query, time.Now(), id)
    if err != nil {
        return fmt.Errorf("failed to delete note: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("note not found")
    }

    return nil
}

// GetNotesByUserID retrieves all notes created by a specific user
func (r *NoteRepo) GetNotesByUserID(userID int) ([]models.Note, error) {
    query := `
        SELECT note_id, user_id, contact_id, lead_id, deal_id, note_text, created_at, updated_at
        FROM notes
        WHERE user_id = $1 AND deleted_at IS NULL
        ORDER BY created_at DESC
    `

    var notes []models.Note
    err := r.db.Select(&notes, query, userID)

    if err != nil {
        return nil, fmt.Errorf("failed to get notes by user ID: %w", err)
    }

    return notes, nil
}

// GetNotesByDealID retrieves all notes for a specific deal
func (r *NoteRepo) GetNotesByDealID(dealID int) ([]models.Note, error) {
    query := `
        SELECT note_id, user_id, contact_id, lead_id, deal_id, note_text, created_at, updated_at
        FROM notes
        WHERE deal_id = $1 AND deleted_at IS NULL
        ORDER BY created_at DESC
    `

    var notes []models.Note
    err := r.db.Select(&notes, query, dealID)

    if err != nil {
        return nil, fmt.Errorf("failed to get notes by deal ID: %w", err)
    }

    return notes, nil
}

// CreateDealNote creates a new note for a deal
func (r *NoteRepo) CreateDealNote(note *models.Note) error {
    if note.DealID == nil {
        return fmt.Errorf("deal ID is required")
    }
    return r.CreateNote(note)
}

// UpdateDealNote updates a deal note
func (r *NoteRepo) UpdateDealNote(note *models.Note) error {
    return r.UpdateNote(note)
}

// DeleteDealNote deletes a deal note
func (r *NoteRepo) DeleteDealNote(id int) error {
    return r.DeleteNote(id)
}