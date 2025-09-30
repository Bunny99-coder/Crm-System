package service

import (
	"errors"
	"fmt"

	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
)

type NoteService struct {
	noteRepo postgres.NoteRepository
}

func NewNoteService(noteRepo postgres.NoteRepository) *NoteService {
	return &NoteService{noteRepo: noteRepo}
}

// CreateNote creates a new note
func (s *NoteService) CreateNote(note *models.Note) error {
	if note.Content == "" {
		return errors.New("note content cannot be empty")
	}
	
	if note.UserID == 0 {
		return errors.New("user ID is required")
	}
	
	return s.noteRepo.CreateNote(note)
}

// GetNoteByID retrieves a note by its ID
func (s *NoteService) GetNoteByID(id int) (*models.Note, error) {
	if id <= 0 {
		return nil, errors.New("invalid note ID")
	}
	
	return s.noteRepo.GetNoteByID(id)
}

// GetNotesByContactID retrieves all notes for a specific contact
func (s *NoteService) GetNotesByContactID(contactID int) ([]models.Note, error) {
	if contactID <= 0 {
		return nil, errors.New("invalid contact ID")
	}
	
	return s.noteRepo.GetNotesByContactID(contactID)
}

// UpdateNote updates an existing note
func (s *NoteService) UpdateNote(note *models.Note) error {
	if note.ID <= 0 {
		return errors.New("invalid note ID")
	}
	
	if note.Content == "" {
		return errors.New("note content cannot be empty")
	}
	
	// Verify the note exists first
	existingNote, err := s.noteRepo.GetNoteByID(note.ID)
	if err != nil {
		return fmt.Errorf("failed to verify note existence: %w", err)
	}
	
	// Ensure user can only update their own notes (optional security check)
	if note.UserID != 0 && existingNote.UserID != note.UserID {
		return errors.New("unauthorized to update this note")
	}
	
	return s.noteRepo.UpdateNote(note)
}

// DeleteNote deletes a note
func (s *NoteService) DeleteNote(id int) error {
	if id <= 0 {
		return errors.New("invalid note ID")
	}
	
	return s.noteRepo.DeleteNote(id)
}

// GetNotesByUserID retrieves all notes created by a specific user
func (s *NoteService) GetNotesByUserID(userID int) ([]models.Note, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	
	return s.noteRepo.GetNotesByUserID(userID)
}







func (s *NoteService) GetNotesByDealID(dealID int) ([]models.Note, error) {
	if dealID <= 0 {
		return nil, errors.New("invalid deal ID")
	}
	return s.noteRepo.GetNotesByDealID(dealID)
}

func (s *NoteService) CreateDealNote(note *models.Note) error {
	if note.Content == "" {
		return errors.New("note content cannot be empty")
	}
	if note.UserID == 0 {
		return errors.New("user ID is required")
	}
	if note.DealID == nil || *note.DealID <= 0 {
		return errors.New("deal ID is required")
	}
	return s.noteRepo.CreateDealNote(note)
}

func (s *NoteService) UpdateDealNote(note *models.Note) error {
	if note.ID <= 0 {
		return errors.New("invalid note ID")
	}
	if note.Content == "" {
		return errors.New("note content cannot be empty")
	}
	return s.noteRepo.UpdateDealNote(note)
}

func (s *NoteService) DeleteDealNote(id int) error {
	if id <= 0 {
		return errors.New("invalid note ID")
	}
	return s.noteRepo.DeleteDealNote(id)
}
