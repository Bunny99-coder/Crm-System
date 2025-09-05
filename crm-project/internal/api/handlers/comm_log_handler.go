// Replace the contents of internal/api/handlers/comm_log_handler.go
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

type CommLogHandler struct {
	service *service.CommLogService
	logger  *slog.Logger
}

func NewCommLogHandler(s *service.CommLogService, logger *slog.Logger) *CommLogHandler {
	return &CommLogHandler{service: s, logger: logger}
}

func (h *CommLogHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newLog models.CommLog
	if err := json.NewDecoder(r.Body).Decode(&newLog); err != nil {
		h.logger.Warn("invalid create comm-log request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newID, err := h.service.CreateLog(ctx, newLog)
	if err != nil {
		h.logger.Error("failed to create comm-log", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.logger.Info("comm-log created successfully", "log_id", newID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": newID})
}

func (h *CommLogHandler) GetLogsForContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contactID, err := strconv.Atoi(chi.URLParam(r, "contactId"))
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}
	logs, err := h.service.GetLogsForContact(ctx, contactID)
	if err != nil {
		h.logger.Error("failed to get comm-logs for contact", "contact_id", contactID, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Debug("retrieved all comm-logs for contact", "contact_id", contactID, "count", len(logs))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *CommLogHandler) GetLogByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "logId"))
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}
	logEntry, err := h.service.GetLogByID(ctx, id)
	if err != nil {
		h.logger.Warn("comm-log not found", "log_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Debug("retrieved comm-log by id", "log_id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntry)
}

func (h *CommLogHandler) UpdateLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "logId"))
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}
	var logToUpdate models.CommLog
	if err := json.NewDecoder(r.Body).Decode(&logToUpdate); err != nil {
		h.logger.Warn("invalid update comm-log request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = h.service.UpdateLog(ctx, id, logToUpdate)
	if err != nil {
		h.logger.Error("failed to update comm-log", "log_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("comm-log updated successfully", "log_id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *CommLogHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "logId"))
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteLog(ctx, id)
	if err != nil {
		h.logger.Error("failed to delete comm-log", "log_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("comm-log deleted successfully", "log_id", id)
	w.WriteHeader(http.StatusNoContent)
}