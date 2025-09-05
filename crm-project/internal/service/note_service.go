// Replace the contents of internal/service/note_service.go
package service

import (
	"context"
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type NoteService struct {
	noteRepo *postgres.NoteRepo
	userRepo *postgres.UserRepo
	logger   *slog.Logger
}

func NewNoteService(nr *postgres.NoteRepo, ur *postgres.UserRepo, logger *slog.Logger) *NoteService {
	return &NoteService{noteRepo: nr, userRepo: ur, logger: logger}
}

func (s *NoteService) CreateNote(ctx context.Context, n models.Note) (int, error) {
	if n.NoteText == "" {
		return 0, errors.New("note text cannot be empty")
	}
	if _, err := s.userRepo.GetByID(ctx, n.UserID); err != nil {
		return 0, fmt.Errorf("invalid user_id: %d", n.UserID)
	}
	if n.NoteDate.IsZero() {
		n.NoteDate = time.Now()
	}
	return s.noteRepo.Create(ctx, n)
}


// in internal/service/note_service.go
func (s *NoteService) GetAllNotes(ctx context.Context) ([]models.Note, error) {
	return s.noteRepo.GetAll(ctx)
}

func (s *NoteService) GetNotesByUser(ctx context.Context, userID int) ([]models.Note, error) {
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("invalid user_id: %d", userID)
	}
	return s.noteRepo.GetAllByUser(ctx, userID)
}

func (s *NoteService) GetNoteByID(ctx context.Context, id int) (*models.Note, error) {
	note, err := s.noteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, fmt.Errorf("note with ID %d not found", id)
	}
	return note, nil
}

func (s *NoteService) UpdateNote(ctx context.Context, id int, n models.Note) error {
	_, err := s.GetNoteByID(ctx, id)
	if err != nil {
		return err
	}
	n.ID = id
	if n.NoteText == "" {
		return errors.New("note text cannot be empty")
	}
	err = s.noteRepo.Update(ctx, n)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("note with ID %d not found during update", id)
		}
		return err
	}
	return nil
}

func (s *NoteService) DeleteNote(ctx context.Context, id int) error {
	err := s.noteRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("note with ID %d not found", id)
		}
		return err
	}
	return nil
}