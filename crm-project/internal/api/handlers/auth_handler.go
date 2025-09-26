// Create new file: internal/api/handlers/auth_handler.go
package handlers

import (
	"crm-project/internal/dto"
	"crm-project/internal/service"
	"encoding/json"
	"net/http"
	"log/slog" 

)

type AuthHandler struct {
	service *service.AuthService
	logger  *slog.Logger // <-- Added

}

func NewAuthHandler(s *service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{service: s, logger: logger}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid login request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	token, roleID, err := h.service.LoginUser(ctx, req.Username, req.Password)
	if err != nil {
		h.logger.Warn("failed login attempt", "username", req.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := LoginResponse{Token: token, RoleID: roleID}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode login response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	h.logger.Info("login response sent", "response", response)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid register request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := h.service.RegisterUser(ctx, &req)
	if err != nil {
		h.logger.Error("failed to register user", "error", err, "username", req.Username)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    // If using JWT: frontend just deletes token, backend can respond OK
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message":"Logged out successfully"}`))
}
