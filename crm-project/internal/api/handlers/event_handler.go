// Replace the contents of internal/api/handlers/event_handler.go
package handlers

import (
	"crm-project/internal/models"
	"crm-project/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	service *service.EventService
	logger  *slog.Logger
}

func NewEventHandler(s *service.EventService, logger *slog.Logger) *EventHandler {
	return &EventHandler{service: s, logger: logger}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newEvent models.Event
	if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
		h.logger.Warn("invalid create event request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newID, err := h.service.CreateEvent(ctx, newEvent)
	if err != nil {
		h.logger.Error("failed to create event", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.logger.Info("event created successfully", "event_id", newID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": newID})
}


// in event_handler.go
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	events, err := h.service.GetAllEvents(ctx)
	if err != nil {
		h.logger.Error("failed to get all events", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.logger.Debug("retrieved all events", "count", len(events))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}



func (h *EventHandler) GetEventsForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	events, err := h.service.GetEventsForUser(ctx, userID)
	if err != nil {
		h.logger.Error("failed to get events for user", "user_id", userID, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Debug("retrieved all events for user", "user_id", userID, "count", len(events))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "eventId"))
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
	event, err := h.service.GetEventByID(ctx, id)
	if err != nil {
		h.logger.Warn("event not found", "event_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Debug("retrieved event by id", "event_id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "eventId"))
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
	var e models.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		h.logger.Warn("invalid update event request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = h.service.UpdateEvent(ctx, id, e)
	if err != nil {
		h.logger.Error("failed to update event", "event_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("event updated successfully", "event_id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "eventId"))
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteEvent(ctx, id)
	if err != nil {
		h.logger.Error("failed to delete event", "event_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("event deleted successfully", "event_id", id)
	w.WriteHeader(http.StatusNoContent)
}