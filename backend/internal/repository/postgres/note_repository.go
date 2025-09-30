package postgres

import "crm-project/internal/models"

// NoteRepository defines the interface for note data operations
type NoteRepository interface {
	CreateNote(note *models.Note) error
	GetNoteByID(id int) (*models.Note, error)
	GetNotesByContactID(contactID int) ([]models.Note, error)
	GetNotesByUserID(userID int) ([]models.Note, error)
	UpdateNote(note *models.Note) error
	DeleteNote(id int) error

	// --- Deal-specific methods ---
	GetNotesByDealID(dealID int) ([]models.Note, error)
	CreateDealNote(note *models.Note) error
	UpdateDealNote(note *models.Note) error
	DeleteDealNote(id int) error
}
