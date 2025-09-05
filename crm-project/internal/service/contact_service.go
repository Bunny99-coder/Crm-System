// Replace the contents of contact_service.go
package service

import (
	"context"
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"database/sql"
	"fmt"
	"log/slog"
)

type ContactService struct {
	repo   *postgres.ContactRepo
	logger *slog.Logger
}

func NewContactService(repo *postgres.ContactRepo, logger *slog.Logger) *ContactService {
	return &ContactService{repo: repo, logger: logger}
}

// All methods are here, ready to use the logger if complex logic is added later.
func (s *ContactService) CreateContact(ctx context.Context, contact models.Contact) (int, error) {
	return s.repo.Create(ctx, contact)
}
func (s *ContactService) GetAllContacts(ctx context.Context) ([]models.Contact, error) {
	return s.repo.GetAll(ctx)
}
func (s *ContactService) GetContactByID(ctx context.Context, id int) (*models.Contact, error) {
	contact, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	if contact == nil { return nil, fmt.Errorf("contact with ID %d not found", id) }
	return contact, nil
}
func (s *ContactService) UpdateContact(ctx context.Context, id int, contact models.Contact) error {
	_, err := s.GetContactByID(ctx, id)
	if err != nil { return err }
	contact.ID = id
	err = s.repo.Update(ctx, contact)
	if err != nil {
		if err == sql.ErrNoRows { return fmt.Errorf("contact with ID %d not found", id) }
		return err
	}
	return nil
}
func (s *ContactService) DeleteContact(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows { return fmt.Errorf("contact with ID %d not found", id) }
		return err
	}
	return nil
}