// Replace the contents of internal/service/event_service.go
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

type EventService struct {
	eventRepo *postgres.EventRepo
	userRepo  *postgres.UserRepo
	logger    *slog.Logger
}

func NewEventService(er *postgres.EventRepo, ur *postgres.UserRepo, logger *slog.Logger) *EventService {
	return &EventService{eventRepo: er, userRepo: ur, logger: logger}
}

func (s *EventService) CreateEvent(ctx context.Context, e models.Event) (int, error) {
	if e.EventName == "" {
		return 0, errors.New("event name is required")
	}
	if e.EndTime.Before(e.StartTime) {
		return 0, errors.New("end time must be after start time")
	}
	if _, err := s.userRepo.GetByID(ctx, e.OrganizerID); err != nil {
		return 0, fmt.Errorf("invalid organizer_id: %d", e.OrganizerID)
	}
	return s.eventRepo.Create(ctx, e)
}

// in event_service.go
func (s *EventService) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	return s.eventRepo.GetAll(ctx)
}

func (s *EventService) GetEventsForUser(ctx context.Context, organizerID int) ([]models.Event, error) {
	if _, err := s.userRepo.GetByID(ctx, organizerID); err != nil {
		return nil, fmt.Errorf("invalid user_id: %d", organizerID)
	}
	return s.eventRepo.GetAllForUser(ctx, organizerID)
}

func (s *EventService) GetEventByID(ctx context.Context, id int) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("event with ID %d not found", id)
	}
	return event, nil
}

func (s *EventService) UpdateEvent(ctx context.Context, id int, e models.Event) error {
	_, err := s.GetEventByID(ctx, id)
	if err != nil {
		return err
	}
	e.ID = id
	if e.EventName == "" {
		return errors.New("event name is required")
	}
	if e.EndTime.Before(e.StartTime) {
		return errors.New("end time must be after start time")
	}

	err = s.eventRepo.Update(ctx, e)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("event with ID %d not found during update", id)
		}
		return err
	}
	return nil
}

func (s *EventService) DeleteEvent(ctx context.Context, id int) error {
	err := s.eventRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("event with ID %d not found", id)
		}
		return err
	}
	return nil
}