package service

import (
	"errors"
	"fmt"

	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
)

type EventService struct {
	eventRepo postgres.EventRepository
}

func NewEventService(eventRepo postgres.EventRepository) *EventService {
	return &EventService{eventRepo: eventRepo}
}

// CreateEvent creates a new event
func (s *EventService) CreateEvent(event *models.Event) error {
	if event.EventName == "" {
		return errors.New("event name cannot be empty")
	}
	if event.OrganizerID == 0 {
		return errors.New("organizer ID is required")
	}
	if event.StartTime.After(event.EndTime) {
		return errors.New("start time must be before end time")
	}

	return s.eventRepo.CreateEvent(event)
}

// GetEventByID retrieves an event by ID
func (s *EventService) GetEventByID(id int) (*models.Event, error) {
	if id <= 0 {
		return nil, errors.New("invalid event ID")
	}

	return s.eventRepo.GetEventByID(id)
}

// GetEventsByDealID retrieves all events for a specific deal
func (s *EventService) GetEventsByDealID(dealID int) ([]models.Event, error) {
	if dealID <= 0 {
		return nil, errors.New("invalid deal ID")
	}

	return s.eventRepo.GetEventsByDealID(dealID)
}

// GetAllEvents retrieves all events
func (s *EventService) GetAllEvents() ([]models.Event, error) {
	return s.eventRepo.GetAllEvents()
}

// UpdateEvent updates an existing event
func (s *EventService) UpdateEvent(event *models.Event) error {
	if event.ID <= 0 {
		return errors.New("invalid event ID")
	}
	if event.EventName == "" {
		return errors.New("event name cannot be empty")
	}
	if event.OrganizerID == 0 {
		return errors.New("organizer ID is required")
	}

	existingEvent, err := s.eventRepo.GetEventByID(event.ID)
	if err != nil {
		return fmt.Errorf("failed to verify event existence: %w", err)
	}

	if existingEvent.OrganizerID != event.OrganizerID {
		return errors.New("unauthorized to update this event")
	}

	return s.eventRepo.UpdateEvent(event)
}

// DeleteEvent soft deletes an event
func (s *EventService) DeleteEvent(id int) error {
	if id <= 0 {
		return errors.New("invalid event ID")
	}

	return s.eventRepo.DeleteEvent(id)
}

// GetEventsForUser retrieves events for a specific user
func (s *EventService) GetEventsForUser(userID int) ([]models.Event, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.eventRepo.GetEventsForUser(userID)
}

// CreateDealEvent creates a new event for a deal (nested route)
func (s *EventService) CreateDealEvent(event *models.Event) error {
	if event.EventName == "" {
		return errors.New("event name cannot be empty")
	}
	if event.OrganizerID == 0 {
		return errors.New("organizer ID is required")
	}
	if event.DealID == nil || *event.DealID <= 0 {
		return errors.New("deal ID is required")
	}
	if event.StartTime.After(event.EndTime) {
		return errors.New("start time must be before end time")
	}

	return s.eventRepo.CreateEvent(event)
}

// UpdateDealEvent updates an event for a deal
func (s *EventService) UpdateDealEvent(event *models.Event) error {
	if event.ID <= 0 {
		return errors.New("invalid event ID")
	}
	if event.EventName == "" {
		return errors.New("event name cannot be empty")
	}

	existingEvent, err := s.eventRepo.GetEventByID(event.ID)
	if err != nil {
		return fmt.Errorf("failed to verify event existence: %w", err)
	}

	if existingEvent.DealID == nil || *existingEvent.DealID != *event.DealID {
		return errors.New("event not found for this deal")
	}

	return s.eventRepo.UpdateEvent(event)
}

// DeleteDealEvent deletes an event for a deal
func (s *EventService) DeleteDealEvent(id int) error {
	if id <= 0 {
		return errors.New("invalid event ID")
	}

	return s.eventRepo.DeleteEvent(id)
}