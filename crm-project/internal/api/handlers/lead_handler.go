// Replace the contents of internal/api/handlers/lead_handler.go
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

type LeadHandler struct {
	service *service.LeadService
	logger  *slog.Logger
}

func NewLeadHandler(s *service.LeadService, logger *slog.Logger) *LeadHandler {
	return &LeadHandler{service: s, logger: logger}
}

func (h *LeadHandler) CreateLead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newLead models.Lead
	if err := json.NewDecoder(r.Body).Decode(&newLead); err != nil {
		h.logger.Warn("invalid create lead request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newID, err := h.service.CreateLead(ctx, newLead)
	if err != nil {
		h.logger.Error("failed to create lead", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.logger.Info("lead created successfully", "lead_id", newID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": newID})
}

func (h *LeadHandler) GetAllLeads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	leads, err := h.service.GetAllLeads(ctx)
	if err != nil {
		h.logger.Error("failed to get all leads", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.logger.Debug("retrieved all leads", "count", len(leads))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leads)
}

func (h *LeadHandler) GetLeadByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid lead ID", http.StatusBadRequest)
		return
	}
	lead, err := h.service.GetLeadByID(ctx, id)
	if err != nil {
		h.logger.Warn("lead not found", "lead_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Debug("retrieved lead by id", "lead_id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lead)
}

func (h *LeadHandler) UpdateLead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid lead ID", http.StatusBadRequest)
		return
	}
	var l models.Lead
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		h.logger.Warn("invalid update lead request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = h.service.UpdateLead(ctx, id, l)
	if err != nil {
		h.logger.Error("failed to update lead", "lead_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("lead updated successfully", "lead_id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *LeadHandler) DeleteLead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid lead ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteLead(ctx, id)
	if err != nil {
		h.logger.Error("failed to delete lead", "lead_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("lead deleted successfully", "lead_id", id)
	w.WriteHeader(http.StatusNoContent)
}