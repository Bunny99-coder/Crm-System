package handlers

import (
	"crm-project/internal/service"
	"encoding/json"
	"log/slog"
	"net/http"
)

type TempHandler struct {
	service *service.TempService
	logger  *slog.Logger
}

func NewTempHandler(s *service.TempService, logger *slog.Logger) *TempHandler {
	return &TempHandler{service: s, logger: logger}
}

func (h *TempHandler) GrantReceptionRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid grant reception role request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.GrantReceptionRole(ctx, req.Username); err != nil {
		h.logger.Error("failed to grant reception role", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("reception role granted successfully", "username", req.Username)
	w.WriteHeader(http.StatusNoContent)
}
