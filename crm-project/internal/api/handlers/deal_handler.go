// Replace the contents of internal/api/handlers/deal_handler.go
package handlers

import (
	"crm-project/internal/models"
	"crm-project/internal/service"
	"crm-project/internal/util"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DealHandler struct {
	service *service.DealService
	logger  *slog.Logger
}

func NewDealHandler(s *service.DealService, logger *slog.Logger) *DealHandler {
	return &DealHandler{service: s, logger: logger}
}

func (h *DealHandler) CreateDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		h.logger.Error("claims not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var newDeal models.Deal
	if err := json.NewDecoder(r.Body).Decode(&newDeal); err != nil {
		h.logger.Warn("invalid create deal request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newDeal.CreatedBy.Int64 = int64(claims.UserID)
	newDeal.CreatedBy.Valid = true

	newID, err := h.service.CreateDeal(ctx, newDeal)
	if err != nil {
		h.logger.Error("failed to create deal", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.logger.Info("deal created successfully", "deal_id", newID, "created_by", claims.UserID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": newID})
}

func (h *DealHandler) GetAllDeals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		h.logger.Error("claims not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	deals, err := h.service.GetAllDeals(ctx, claims.UserID, claims.RoleID)
	if err != nil {
		h.logger.Error("failed to get all deals", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.logger.Debug("retrieved all deals", "count", len(deals), "user_id", claims.UserID, "role_id", claims.RoleID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deals)
}

func (h *DealHandler) GetDealByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		h.logger.Error("claims not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid deal ID", http.StatusBadRequest)
		return
	}
	deal, err := h.service.GetDealByID(ctx, id, claims.UserID, claims.RoleID)
	if err != nil {
		h.logger.Warn("deal not found or unauthorized", "deal_id", id, "user_id", claims.UserID, "role_id", claims.RoleID, "error", err)
		if err.Error() == "unauthorized" {
			http.Error(w, "Unauthorized to view this deal", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		return
	}
	h.logger.Debug("retrieved deal by id", "deal_id", id, "user_id", claims.UserID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deal)
}

func (h *DealHandler) UpdateDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		h.logger.Error("claims not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid deal ID", http.StatusBadRequest)
		return
	}
	var d models.Deal
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.logger.Warn("invalid update deal request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = h.service.UpdateDeal(ctx, id, d, claims.UserID, claims.RoleID)
	if err != nil {
		h.logger.Error("failed to update deal", "deal_id", id, "user_id", claims.UserID, "role_id", claims.RoleID, "error", err)
		if err.Error() == "unauthorized" {
			http.Error(w, "Unauthorized to update this deal", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	h.logger.Info("deal updated successfully", "deal_id", id, "user_id", claims.UserID)
	w.WriteHeader(http.StatusNoContent)
}

func (h *DealHandler) DeleteDeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		h.logger.Error("claims not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid deal ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteDeal(ctx, id, claims.UserID, claims.RoleID)
	if err != nil {
		h.logger.Error("failed to delete deal", "deal_id", id, "user_id", claims.UserID, "role_id", claims.RoleID, "error", err)
		if err.Error() == "unauthorized" {
			http.Error(w, "Unauthorized to delete this deal", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	h.logger.Info("deal deleted successfully", "deal_id", id, "user_id", claims.UserID)
	w.WriteHeader(http.StatusNoContent)
}