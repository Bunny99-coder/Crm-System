// Replace the contents of internal/api/handlers/property_handler.go
package handlers

import (
	"crm-project/internal/models"
	"crm-project/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings" 
	"github.com/go-chi/chi/v5"
)

type PropertyHandler struct {
	service *service.PropertyService
	logger  *slog.Logger
}

func NewPropertyHandler(s *service.PropertyService, logger *slog.Logger) *PropertyHandler {
	return &PropertyHandler{service: s, logger: logger}
}

func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newProperty models.Property
	if err := json.NewDecoder(r.Body).Decode(&newProperty); err != nil {
		h.logger.Warn("invalid create property request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newID, err := h.service.CreateProperty(ctx, newProperty)
	if err != nil {
		h.logger.Error("failed to create property", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest) // Assuming validation errors
		return
	}
	h.logger.Info("property created successfully", "property_id", newID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": newID})
}

func (h *PropertyHandler) GetAllProperties(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	properties, err := h.service.GetAllProperties(ctx)
	if err != nil {
		h.logger.Error("failed to get all properties", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.logger.Debug("retrieved all properties", "count", len(properties))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(properties)
}

func (h *PropertyHandler) GetPropertyByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}
	property, err := h.service.GetPropertyByID(ctx, id)
	if err != nil {
		h.logger.Warn("property not found", "property_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Debug("retrieved property by id", "property_id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(property)
}

// Replace the UpdateProperty function in property_handler.go with this:
func (h *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "propertyId") // <-- THE FIX
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid property ID format", http.StatusBadRequest)
		return
	}
	var p models.Property
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.logger.Warn("invalid update property request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = h.service.UpdateProperty(ctx, id, p)
	if err != nil {
		h.logger.Error("failed to update property", "property_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("property updated successfully", "property_id", id)
	w.WriteHeader(http.StatusNoContent)
}

// Replace the DeleteProperty function in property_handler.go with this:
func (h *PropertyHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "propertyId") // <-- THE FIX
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid property ID format", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteProperty(ctx, id)
	if err != nil {
		// Also check for foreign key constraint errors
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			h.logger.Warn("attempted to delete property with dependent leads/deals", "property_id", id)
			http.Error(w, "Cannot delete property: it is linked to existing leads or deals.", http.StatusConflict) // 409 Conflict is a good code here
			return
		}
		h.logger.Error("failed to delete property", "property_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("property deleted successfully", "property_id", id)
	w.WriteHeader(http.StatusNoContent)
}