// Replace the contents of internal/api/handlers/user_handler.go
package handlers

import (
	"crm-project/internal/dto"
	"crm-project/internal/service"
	"crm-project/internal/util"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service *service.UserService
	logger  *slog.Logger
}

func NewUserHandler(s *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{service: s, logger: logger}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid create user request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newID, err := h.service.CreateUser(ctx, req)
	if err != nil {
		if _, ok := err.(*util.ValidationError); ok {
			h.logger.Warn("create user validation failed", "error", err, "request", req)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			h.logger.Error("failed to create user", "error", err)
			// A real app might check for specific DB errors, like unique constraint violation
			http.Error(w, "Could not create user", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("user created successfully", "user_id", newID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": newID})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.service.GetAllUsers(ctx)
	if err != nil {
		h.logger.Error("failed to get all users", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.logger.Debug("retrieved all users", "count", len(users))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		h.logger.Warn("user not found", "user_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Debug("retrieved user by id", "user_id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid update user request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateUser(ctx, id, req)
	if err != nil {
		if _, ok := err.(*util.ValidationError); ok {
			h.logger.Warn("update user validation failed", "error", err, "user_id", id)
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			h.logger.Error("failed to update user", "error", err, "user_id", id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.logger.Info("user updated successfully", "user_id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteUser(ctx, id)
	if err != nil {
		h.logger.Error("failed to delete user", "error", err, "user_id", id)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("user deleted successfully", "user_id", id)
	w.WriteHeader(http.StatusNoContent)
}