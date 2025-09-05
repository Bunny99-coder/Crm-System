// Replace the contents of internal/service/comm_log_service.go
package service

import (
	"context"
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type CommLogService struct {
	commRepo    *postgres.CommLogRepo
	contactRepo *postgres.ContactRepo
	userRepo    *postgres.UserRepo
	logger      *slog.Logger
}

func NewCommLogService(clr *postgres.CommLogRepo, cr *postgres.ContactRepo, ur *postgres.UserRepo, logger *slog.Logger) *CommLogService {
	return &CommLogService{
		commRepo:    clr,
		contactRepo: cr,
		userRepo:    ur,
		logger:      logger,
	}
}

func (s *CommLogService) CreateLog(ctx context.Context, cl models.CommLog) (int, error) {
	if cl.InteractionType == "" {
		return 0, errors.New("interaction type is required")
	}
	if _, err := s.contactRepo.GetByID(ctx, cl.ContactID); err != nil {
		return 0, fmt.Errorf("invalid contact_id: %d", cl.ContactID)
	}
	if _, err := s.userRepo.GetByID(ctx, cl.UserID); err != nil {
		return 0, fmt.Errorf("invalid user_id: %d", cl.UserID)
	}
	return s.commRepo.Create(ctx, cl)
}

func (s *CommLogService) GetLogsForContact(ctx context.Context, contactID int) ([]models.CommLog, error) {
	if _, err := s.contactRepo.GetByID(ctx, contactID); err != nil {
		return nil, fmt.Errorf("invalid contact_id: %d", contactID)
	}
	return s.commRepo.GetAllForContact(ctx, contactID)
}

func (s *CommLogService) GetLogByID(ctx context.Context, id int) (*models.CommLog, error) {
	log, err := s.commRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if log == nil {
		return nil, fmt.Errorf("log with ID %d not found", id)
	}
	return log, nil
}

func (s *CommLogService) UpdateLog(ctx context.Context, id int, cl models.CommLog) error {
	_, err := s.GetLogByID(ctx, id)
	if err != nil {
		return err
	}
	cl.ID = id
	if cl.InteractionType == "" {
		return errors.New("interaction type cannot be empty")
	}
	err = s.commRepo.Update(ctx, cl)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("log with ID %d not found during update", id)
		}
		return err
	}
	return nil
}

func (s *CommLogService) DeleteLog(ctx context.Context, id int) error {
	err := s.commRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("log with ID %d not found", id)
		}
		return err
	}
	return nil
}