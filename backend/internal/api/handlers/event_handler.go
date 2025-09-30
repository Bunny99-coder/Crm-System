package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"crm-project/internal/models"
	"crm-project/internal/service"
	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	eventService *service.EventService
}

func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

// CreateEventRequest represents the request body for creating an event
type CreateEventRequest struct {
	EventName        string     `json:"event_name"`
	EventDescription *string    `json:"event_description,omitempty"`
	StartTime        time.Time  `json:"start_time"`
	EndTime          time.Time  `json:"end_time"`
	Location         *string    `json:"location,omitempty"`
	LeadID           *int       `json:"lead_id,omitempty"`
	DealID           *int       `json:"deal_id,omitempty"`
}

// UpdateEventRequest represents the request body for updating an event
type UpdateEventRequest struct {
	EventName        string     `json:"event_name"`
	EventDescription *string    `json:"event_description,omitempty"`
	StartTime        time.Time  `json:"start_time"`
	EndTime          time.Time  `json:"end_time"`
	Location         *string    `json:"location,omitempty"`
}

// EventResponse represents the response structure for events
type EventResponse struct {
	ID               int        `json:"id"`
	EventName        string     `json:"event_name"`
	EventDescription *string    `json:"event_description,omitempty"`
	StartTime        time.Time  `json:"start_time"`
	EndTime          time.Time  `json:"end_time"`
	Location         *string    `json:"location,omitempty"`
	OrganizerID      int        `json:"organizer_id"`
	LeadID           *int       `json:"lead_id,omitempty"`
	DealID           *int       `json:"deal_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

// Helper function to convert models.Event to EventResponse
func convertEventToResponse(event *models.Event) EventResponse {
	response := EventResponse{
		ID:               event.ID,
		EventName:        event.EventName,
		EventDescription: event.EventDescription,
		StartTime:        event.StartTime,
		EndTime:          event.EndTime,
		Location:         event.Location,
		OrganizerID:      event.OrganizerID,
		LeadID:           event.LeadID,
		DealID:           event.DealID,
		CreatedAt:        event.CreatedAt,
		UpdatedAt:        event.UpdatedAt,
	}
	return response
}

// GetAllEvents handles GET /api/v1/events
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	slog.Info("GetAllEvents called", "method", r.Method, "url", r.URL.Path)

	events, err := h.eventService.GetAllEvents()
	if err != nil {
		slog.Error("Failed to get all events", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get events: "+err.Error())
		return
	}

	eventResponses := make([]EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = convertEventToResponse(&event)
	}

	slog.Info("Successfully retrieved all events", "count", len(events))
	respondWithJSON(w, http.StatusOK, eventResponses)
}

// CreateEvent handles POST /api/v1/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	slog.Info("CreateEvent called", "method", r.Method, "url", r.URL.Path)
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	event := &models.Event{
		EventName:        req.EventName,
		EventDescription: req.EventDescription,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Location:         req.Location,
		OrganizerID:      userID,
		LeadID:           req.LeadID,
		DealID:           req.DealID,
		CreatedAt:        time.Now(),
	}

	if err := h.eventService.CreateEvent(event); err != nil {
		slog.Error("Failed to create event", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create event: "+err.Error())
		return
	}

	slog.Info("Successfully created event", "eventID", event.ID)
	respondWithJSON(w, http.StatusCreated, convertEventToResponse(event))
}

// GetEventByID handles GET /api/v1/events/{eventId}
func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("GetEventByID called", "method", r.Method, "url", r.URL.Path)
	eventID, err := strconv.Atoi(chi.URLParam(r, "eventId"))
	if err != nil {
		slog.Error("Invalid event ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	event, err := h.eventService.GetEventByID(eventID)
	if err != nil {
		slog.Error("Failed to get event", "eventID", eventID, "error", err)
		respondWithError(w, http.StatusNotFound, "Event not found: "+err.Error())
		return
	}

	slog.Info("Successfully retrieved event", "eventID", eventID)
	respondWithJSON(w, http.StatusOK, convertEventToResponse(event))
}

// UpdateEvent handles PUT /api/v1/events/{eventId}
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	slog.Info("UpdateEvent called", "method", r.Method, "url", r.URL.Path)
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	eventID, err := strconv.Atoi(chi.URLParam(r, "eventId"))
	if err != nil {
		slog.Error("Invalid event ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	event := &models.Event{
		ID:               eventID,
		EventName:        req.EventName,
		EventDescription: req.EventDescription,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Location:         req.Location,
		OrganizerID:      userID,
	}

	if err := h.eventService.UpdateEvent(event); err != nil {
		slog.Error("Failed to update event", "eventID", eventID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update event: "+err.Error())
		return
	}

	updatedEvent, err := h.eventService.GetEventByID(eventID)
	if err != nil {
		slog.Error("Failed to fetch updated event", "eventID", eventID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get updated event: "+err.Error())
		return
	}

	slog.Info("Successfully updated event", "eventID", eventID)
	respondWithJSON(w, http.StatusOK, convertEventToResponse(updatedEvent))
}

// DeleteEvent handles DELETE /api/v1/events/{eventId}
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	slog.Info("DeleteEvent called", "method", r.Method, "url", r.URL.Path)
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	eventID, err := strconv.Atoi(chi.URLParam(r, "eventId"))
	if err != nil {
		slog.Error("Invalid event ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	existingEvent, err := h.eventService.GetEventByID(eventID)
	if err != nil {
		slog.Error("Event not found", "eventID", eventID, "error", err)
		respondWithError(w, http.StatusNotFound, "Event not found: "+err.Error())
		return
	}

	if existingEvent.OrganizerID != userID {
		slog.Warn("Unauthorized to delete event", "userID", userID, "eventOrganizerID", existingEvent.OrganizerID)
		respondWithError(w, http.StatusForbidden, "You can only delete your own events")
		return
	}

	if err := h.eventService.DeleteEvent(eventID); err != nil {
		slog.Error("Failed to delete event", "eventID", eventID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete event: "+err.Error())
		return
	}

	slog.Info("Successfully deleted event", "eventID", eventID)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Event deleted successfully"})
}

// GetEventsForUser handles GET /api/v1/users/{userId}/events
func (h *EventHandler) GetEventsForUser(w http.ResponseWriter, r *http.Request) {
	slog.Info("GetEventsForUser called", "method", r.Method, "url", r.URL.Path)
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		slog.Error("Invalid user ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	events, err := h.eventService.GetEventsForUser(userID)
	if err != nil {
		slog.Error("Failed to get user events", "userID", userID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get user events: "+err.Error())
		return
	}

	eventResponses := make([]EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = convertEventToResponse(&event)
	}

	slog.Info("Successfully retrieved user events", "userID", userID, "count", len(events))
	respondWithJSON(w, http.StatusOK, eventResponses)
}

// GetDealEvents handles GET /api/v1/deals/{dealId}/events
func (h *EventHandler) GetDealEvents(w http.ResponseWriter, r *http.Request) {
	slog.Info("GetDealEvents called", "method", r.Method, "url", r.URL.Path)
	dealID, err := strconv.Atoi(chi.URLParam(r, "dealId"))
	if err != nil {
		slog.Error("Invalid deal ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	events, err := h.eventService.GetEventsByDealID(dealID)
	if err != nil {
		slog.Error("Failed to get deal events", "dealID", dealID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get deal events: "+err.Error())
		return
	}

	eventResponses := make([]EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = convertEventToResponse(&event)
	}

	slog.Info("Successfully retrieved deal events", "dealID", dealID, "count", len(events))
	respondWithJSON(w, http.StatusOK, eventResponses)
}

// CreateDealEvent handles POST /api/v1/deals/{dealId}/events
func (h *EventHandler) CreateDealEvent(w http.ResponseWriter, r *http.Request) {
	slog.Info("CreateDealEvent called", "method", r.Method, "url", r.URL.Path)
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	dealID, err := strconv.Atoi(chi.URLParam(r, "dealId"))
	if err != nil {
		slog.Error("Invalid deal ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	event := &models.Event{
		EventName:        req.EventName,
		EventDescription: req.EventDescription,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Location:         req.Location,
		OrganizerID:      userID,
		DealID:           &dealID,
		CreatedAt:        time.Now(),
	}

	if err := h.eventService.CreateDealEvent(event); err != nil {
		slog.Error("Failed to create deal event", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create deal event: "+err.Error())
		return
	}

	slog.Info("Successfully created deal event", "eventID", event.ID, "dealID", dealID)
	respondWithJSON(w, http.StatusCreated, convertEventToResponse(event))
}

// UpdateDealEvent handles PUT /api/v1/deals/{dealId}/events/{eventId}
func (h *EventHandler) UpdateDealEvent(w http.ResponseWriter, r *http.Request) {
	slog.Info("UpdateDealEvent called", "method", r.Method, "url", r.URL.Path)
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	dealID, err := strconv.Atoi(chi.URLParam(r, "dealId"))
	if err != nil {
		slog.Error("Invalid deal ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	eventID, err := strconv.Atoi(chi.URLParam(r, "eventId"))
	if err != nil {
		slog.Error("Invalid event ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	event := &models.Event{
		ID:               eventID,
		EventName:        req.EventName,
		EventDescription: req.EventDescription,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Location:         req.Location,
		OrganizerID:      userID,
		DealID:           &dealID,
	}

	if err := h.eventService.UpdateDealEvent(event); err != nil {
		slog.Error("Failed to update deal event", "eventID", eventID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update deal event: "+err.Error())
		return
	}

	updatedEvent, err := h.eventService.GetEventByID(eventID)
	if err != nil {
		slog.Error("Failed to fetch updated event", "eventID", eventID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get updated event: "+err.Error())
		return
	}

	slog.Info("Successfully updated deal event", "eventID", eventID, "dealID", dealID)
	respondWithJSON(w, http.StatusOK, convertEventToResponse(updatedEvent))
}

// DeleteDealEvent handles DELETE /api/v1/deals/{dealId}/events/{eventId}
func (h *EventHandler) DeleteDealEvent(w http.ResponseWriter, r *http.Request) {
	slog.Info("DeleteDealEvent called", "method", r.Method, "url", r.URL.Path)
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	dealID, err := strconv.Atoi(chi.URLParam(r, "dealId"))
	if err != nil {
		slog.Error("Invalid deal ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	eventID, err := strconv.Atoi(chi.URLParam(r, "eventId"))
	if err != nil {
		slog.Error("Invalid event ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	existingEvent, err := h.eventService.GetEventByID(eventID)
	if err != nil {
		slog.Error("Event not found", "eventID", eventID, "error", err)
		respondWithError(w, http.StatusNotFound, "Event not found: "+err.Error())
		return
	}

	if existingEvent.DealID == nil || *existingEvent.DealID != dealID {
		slog.Warn("Event does not belong to deal", "eventID", eventID, "dealID", dealID)
		respondWithError(w, http.StatusNotFound, "Event not found for this deal")
		return
	}

	if existingEvent.OrganizerID != userID {
		slog.Warn("Unauthorized to delete event", "userID", userID, "eventOrganizerID", existingEvent.OrganizerID)
		respondWithError(w, http.StatusForbidden, "You can only delete your own events")
		return
	}

	if err := h.eventService.DeleteEvent(eventID); err != nil {
		slog.Error("Failed to delete deal event", "eventID", eventID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete deal event: "+err.Error())
		return
	}

	slog.Info("Successfully deleted deal event", "eventID", eventID, "dealID", dealID)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Deal event deleted successfully"})
}