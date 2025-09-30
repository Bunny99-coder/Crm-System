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

type CommLogHandler struct {
    commLogService *service.CommLogService
}

func NewCommLogHandler(commLogService *service.CommLogService) *CommLogHandler {
    return &CommLogHandler{commLogService: commLogService}
}

// CreateCommLogRequest represents the request body for creating a communication log
type CreateCommLogRequest struct {
    ContactID       *int      `json:"contact_id,omitempty"`
    LeadID          *int      `json:"lead_id,omitempty"`
    DealID          *int      `json:"deal_id,omitempty"`
    InteractionDate time.Time `json:"interaction_date"`
    InteractionType string    `json:"interaction_type"`
    Notes           *string   `json:"notes,omitempty"`
}

// UpdateCommLogRequest represents the request body for updating a communication log
type UpdateCommLogRequest struct {
    ContactID       *int      `json:"contact_id,omitempty"`
    InteractionDate time.Time `json:"interaction_date"`
    InteractionType string    `json:"interaction_type"`
    Notes           *string   `json:"notes,omitempty"`
}

// CommLogResponse represents the response structure for communication logs
type CommLogResponse struct {
    ID              int        `json:"id"`
    ContactID       *int       `json:"contact_id,omitempty"`
    UserID          int        `json:"user_id"`
    LeadID          *int       `json:"lead_id,omitempty"`
    DealID          *int       `json:"deal_id,omitempty"`
    InteractionDate time.Time  `json:"interaction_date"`
    InteractionType string     `json:"interaction_type"`
    Notes           *string    `json:"notes,omitempty"`
    CreatedAt       *time.Time `json:"created_at,omitempty"`
}

// Helper function to convert models.CommLog to CommLogResponse
func convertCommLogToResponse(log *models.CommLog) CommLogResponse {
    response := CommLogResponse{
        ID:              log.ID,
        ContactID:       log.ContactID,
        UserID:          log.UserID,
        LeadID:          log.LeadID,
        DealID:          log.DealID,
        InteractionDate: log.InteractionDate,
        InteractionType: log.InteractionType,
        Notes:           log.Notes,
        CreatedAt:       log.CreatedAt,
    }
    return response
}

// GetAllCommLogs handles GET /api/v1/comm-logs
func (h *CommLogHandler) GetAllCommLogs(w http.ResponseWriter, r *http.Request) {
    slog.Info("GetAllCommLogs called", "method", r.Method, "url", r.URL.Path)

    logs, err := h.commLogService.GetAllCommLogs()
    if err != nil {
        slog.Error("Failed to get all communication logs", "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get communication logs: "+err.Error())
        return
    }

    logResponses := make([]CommLogResponse, len(logs))
    for i, log := range logs {
        logResponses[i] = convertCommLogToResponse(&log)
    }

    slog.Info("Successfully retrieved all communication logs", "count", len(logs))
    respondWithJSON(w, http.StatusOK, logResponses)
}

// CreateCommLog handles POST /api/v1/comm-logs
func (h *CommLogHandler) CreateCommLog(w http.ResponseWriter, r *http.Request) {
    slog.Info("CreateCommLog called", "method", r.Method, "url", r.URL.Path)
    userID, err := getUserIDFromContext(r.Context())
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    var req CreateCommLogRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("Failed to decode request body", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    log := &models.CommLog{
        ContactID:       req.ContactID,
        UserID:          userID,
        LeadID:          req.LeadID,
        DealID:          req.DealID,
        InteractionDate: req.InteractionDate,
        InteractionType: req.InteractionType,
        Notes:           req.Notes,
        CreatedAt:       &time.Time{},
    }

    if err := h.commLogService.CreateCommLog(log); err != nil {
        slog.Error("Failed to create communication log", "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to create communication log: "+err.Error())
        return
    }

    slog.Info("Successfully created communication log", "logID", log.ID)
    respondWithJSON(w, http.StatusCreated, convertCommLogToResponse(log))
}

// GetCommLogByID handles GET /api/v1/comm-logs/{logId}
func (h *CommLogHandler) GetCommLogByID(w http.ResponseWriter, r *http.Request) {
    slog.Info("GetCommLogByID called", "method", r.Method, "url", r.URL.Path)
    logID, err := strconv.Atoi(chi.URLParam(r, "logId"))
    if err != nil {
        slog.Error("Invalid communication log ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid communication log ID")
        return
    }

    log, err := h.commLogService.GetCommLogByID(logID)
    if err != nil {
        slog.Error("Failed to get communication log", "logID", logID, "error", err)
        respondWithError(w, http.StatusNotFound, "Communication log not found: "+err.Error())
        return
    }

    slog.Info("Successfully retrieved communication log", "logID", logID)
    respondWithJSON(w, http.StatusOK, convertCommLogToResponse(log))
}

// UpdateCommLog handles PUT /api/v1/comm-logs/{logId}
func (h *CommLogHandler) UpdateCommLog(w http.ResponseWriter, r *http.Request) {
    slog.Info("UpdateCommLog called", "method", r.Method, "url", r.URL.Path)
    userID, err := getUserIDFromContext(r.Context())
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    logID, err := strconv.Atoi(chi.URLParam(r, "logId"))
    if err != nil {
        slog.Error("Invalid communication log ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid communication log ID")
        return
    }

    var req UpdateCommLogRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("Failed to decode request body", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    log := &models.CommLog{
        ID:              logID,
        ContactID:       req.ContactID,
        UserID:          userID,
        InteractionDate: req.InteractionDate,
        InteractionType: req.InteractionType,
        Notes:           req.Notes,
    }

    if err := h.commLogService.UpdateCommLog(log); err != nil {
        slog.Error("Failed to update communication log", "logID", logID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to update communication log: "+err.Error())
        return
    }

    updatedLog, err := h.commLogService.GetCommLogByID(logID)
    if err != nil {
        slog.Error("Failed to fetch updated communication log", "logID", logID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get updated communication log: "+err.Error())
        return
    }

    slog.Info("Successfully updated communication log", "logID", logID)
    respondWithJSON(w, http.StatusOK, convertCommLogToResponse(updatedLog))
}

// DeleteCommLog handles DELETE /api/v1/comm-logs/{logId}
func (h *CommLogHandler) DeleteCommLog(w http.ResponseWriter, r *http.Request) {
    slog.Info("DeleteCommLog called", "method", r.Method, "url", r.URL.Path)
    userID, err := getUserIDFromContext(r.Context())
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    logID, err := strconv.Atoi(chi.URLParam(r, "logId"))
    if err != nil {
        slog.Error("Invalid communication log ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid communication log ID")
        return
    }

    existingLog, err := h.commLogService.GetCommLogByID(logID)
    if err != nil {
        slog.Error("Communication log not found", "logID", logID, "error", err)
        respondWithError(w, http.StatusNotFound, "Communication log not found: "+err.Error())
        return
    }

    if existingLog.UserID != userID {
        slog.Warn("Unauthorized to delete communication log", "userID", userID, "logUserID", existingLog.UserID)
        respondWithError(w, http.StatusForbidden, "You can only delete your own communication logs")
        return
    }

    if err := h.commLogService.DeleteCommLog(logID); err != nil {
        slog.Error("Failed to delete communication log", "logID", logID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to delete communication log: "+err.Error())
        return
    }

    slog.Info("Successfully deleted communication log", "logID", logID)
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Communication log deleted successfully"})
}

// GetLogsForContact handles GET /api/v1/contacts/{contactId}/comm-logs
func (h *CommLogHandler) GetLogsForContact(w http.ResponseWriter, r *http.Request) {
    slog.Info("GetLogsForContact called", "method", r.Method, "url", r.URL.Path)
    contactID, err := strconv.Atoi(chi.URLParam(r, "contactId"))
    if err != nil {
        slog.Error("Invalid contact ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
        return
    }

    logs, err := h.commLogService.GetCommLogsByContactID(contactID)
    if err != nil {
        slog.Error("Failed to get communication logs for contact", "contactID", contactID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get communication logs: "+err.Error())
        return
    }

    logResponses := make([]CommLogResponse, len(logs))
    for i, log := range logs {
        logResponses[i] = convertCommLogToResponse(&log)
    }

    slog.Info("Successfully retrieved communication logs for contact", "contactID", contactID, "count", len(logs))
    respondWithJSON(w, http.StatusOK, logResponses)
}

// CreateContactCommLog handles POST /api/v1/contacts/{contactId}/comm-logs
func (h *CommLogHandler) CreateContactCommLog(w http.ResponseWriter, r *http.Request) {
    slog.Info("CreateContactCommLog called", "method", r.Method, "url", r.URL.Path)
    userID, err := getUserIDFromContext(r.Context())
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    contactID, err := strconv.Atoi(chi.URLParam(r, "contactId"))
    if err != nil {
        slog.Error("Invalid contact ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
        return
    }

    var req CreateCommLogRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("Failed to decode request body", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    log := &models.CommLog{
        ContactID:       &contactID,
        UserID:          userID,
        LeadID:          req.LeadID,
        DealID:          req.DealID,
        InteractionDate: req.InteractionDate,
        InteractionType: req.InteractionType,
        Notes:           req.Notes,
        CreatedAt:       &time.Time{},
    }

    if err := h.commLogService.CreateContactCommLog(log); err != nil {
        slog.Error("Failed to create contact communication log", "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to create communication log: "+err.Error())
        return
    }

    slog.Info("Successfully created contact communication log", "logID", log.ID, "contactID", contactID)
    respondWithJSON(w, http.StatusCreated, convertCommLogToResponse(log))
}

// UpdateContactCommLog handles PUT /api/v1/contacts/{contactId}/comm-logs/{logId}
func (h *CommLogHandler) UpdateContactCommLog(w http.ResponseWriter, r *http.Request) {
    slog.Info("UpdateContactCommLog called", "method", r.Method, "url", r.URL.Path)
    userID, err := getUserIDFromContext(r.Context())
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    contactID, err := strconv.Atoi(chi.URLParam(r, "contactId"))
    if err != nil {
        slog.Error("Invalid contact ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
        return
    }

    logID, err := strconv.Atoi(chi.URLParam(r, "logId"))
    if err != nil {
        slog.Error("Invalid communication log ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid communication log ID")
        return
    }

    var req UpdateCommLogRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("Failed to decode request body", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    log := &models.CommLog{
        ID:              logID,
        ContactID:       &contactID,
        UserID:          userID,
        InteractionDate: req.InteractionDate,
        InteractionType: req.InteractionType,
        Notes:           req.Notes,
    }

    if err := h.commLogService.UpdateContactCommLog(log); err != nil {
        slog.Error("Failed to update contact communication log", "logID", logID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to update communication log: "+err.Error())
        return
    }

    updatedLog, err := h.commLogService.GetCommLogByID(logID)
    if err != nil {
        slog.Error("Failed to fetch updated communication log", "logID", logID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get updated communication log: "+err.Error())
        return
    }

    slog.Info("Successfully updated contact communication log", "logID", logID, "contactID", contactID)
    respondWithJSON(w, http.StatusOK, convertCommLogToResponse(updatedLog))
}

// DeleteContactCommLog handles DELETE /api/v1/contacts/{contactId}/comm-logs/{logId}
func (h *CommLogHandler) DeleteContactCommLog(w http.ResponseWriter, r *http.Request) {
    slog.Info("DeleteContactCommLog called", "method", r.Method, "url", r.URL.Path)
    userID, err := getUserIDFromContext(r.Context())
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Authentication required")
        return
    }

    contactID, err := strconv.Atoi(chi.URLParam(r, "contactId"))
    if err != nil {
        slog.Error("Invalid contact ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
        return
    }

    logID, err := strconv.Atoi(chi.URLParam(r, "logId"))
    if err != nil {
        slog.Error("Invalid communication log ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid communication log ID")
        return
    }

    existingLog, err := h.commLogService.GetCommLogByID(logID)
    if err != nil {
        slog.Error("Communication log not found", "logID", logID, "error", err)
        respondWithError(w, http.StatusNotFound, "Communication log not found: "+err.Error())
        return
    }

    if existingLog.ContactID == nil || (existingLog.ContactID != nil && *existingLog.ContactID != contactID) {
        slog.Warn("Communication log does not belong to contact", "logID", logID, "contactID", contactID)
        respondWithError(w, http.StatusNotFound, "Communication log not found for this contact")
        return
    }

    if existingLog.UserID != userID {
        slog.Warn("Unauthorized to delete communication log", "userID", userID, "logUserID", existingLog.UserID)
        respondWithError(w, http.StatusForbidden, "You can only delete your own communication logs")
        return
    }

    if err := h.commLogService.DeleteContactCommLog(logID); err != nil {
        slog.Error("Failed to delete contact communication log", "logID", logID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to delete communication log: "+err.Error())
        return
    }

    slog.Info("Successfully deleted contact communication log", "logID", logID, "contactID", contactID)
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Communication log deleted successfully"})
}

// GetDealCommLogs handles GET /api/v1/deals/{dealId}/comm-logs
func (h *CommLogHandler) GetDealCommLogs(w http.ResponseWriter, r *http.Request) {
    slog.Info("GetDealCommLogs called", "method", r.Method, "url", r.URL.Path)
    dealID, err := strconv.Atoi(chi.URLParam(r, "dealId"))
    if err != nil {
        slog.Error("Invalid deal ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
        return
    }

    logs, err := h.commLogService.GetCommLogsByDealID(dealID)
    if err != nil {
        slog.Error("Failed to get deal communication logs", "dealID", dealID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get deal communication logs: "+err.Error())
        return
    }

    logResponses := make([]CommLogResponse, len(logs))
    for i, log := range logs {
        logResponses[i] = convertCommLogToResponse(&log)
    }

    slog.Info("Successfully retrieved deal communication logs", "dealID", dealID, "count", len(logs))
    respondWithJSON(w, http.StatusOK, logResponses)
}

// CreateDealCommLog handles POST /api/v1/deals/{dealId}/comm-logs
func (h *CommLogHandler) CreateDealCommLog(w http.ResponseWriter, r *http.Request) {
    slog.Info("CreateDealCommLog called", "method", r.Method, "url", r.URL.Path)
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

    var req CreateCommLogRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("Failed to decode request body", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    log := &models.CommLog{
        ContactID:       req.ContactID,
        UserID:          userID,
        DealID:          &dealID,
        InteractionDate: req.InteractionDate,
        InteractionType: req.InteractionType,
        Notes:           req.Notes,
        CreatedAt:       &time.Time{},
    }

    if err := h.commLogService.CreateDealCommLog(log); err != nil {
        slog.Error("Failed to create deal communication log", "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to create deal communication log: "+err.Error())
        return
    }

    slog.Info("Successfully created deal communication log", "logID", log.ID, "dealID", dealID)
    respondWithJSON(w, http.StatusCreated, convertCommLogToResponse(log))
}

// UpdateDealCommLog handles PUT /api/v1/deals/{dealId}/comm-logs/{logId}
func (h *CommLogHandler) UpdateDealCommLog(w http.ResponseWriter, r *http.Request) {
    slog.Info("UpdateDealCommLog called", "method", r.Method, "url", r.URL.Path)
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

    logID, err := strconv.Atoi(chi.URLParam(r, "logId"))
    if err != nil {
        slog.Error("Invalid communication log ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid communication log ID")
        return
    }

    var req UpdateCommLogRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        slog.Error("Failed to decode request body", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    log := &models.CommLog{
        ID:              logID,
        ContactID:       req.ContactID,
        UserID:          userID,
        DealID:          &dealID,
        InteractionDate: req.InteractionDate,
        InteractionType: req.InteractionType,
        Notes:           req.Notes,
    }

    if err := h.commLogService.UpdateDealCommLog(log); err != nil {
        slog.Error("Failed to update deal communication log", "logID", logID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to update deal communication log: "+err.Error())
        return
    }

    updatedLog, err := h.commLogService.GetCommLogByID(logID)
    if err != nil {
        slog.Error("Failed to fetch updated communication log", "logID", logID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to get updated communication log: "+err.Error())
        return
    }

    slog.Info("Successfully updated deal communication log", "logID", logID, "dealID", dealID)
    respondWithJSON(w, http.StatusOK, convertCommLogToResponse(updatedLog))
}

// DeleteDealCommLog handles DELETE /api/v1/deals/{dealId}/comm-logs/{logId}
func (h *CommLogHandler) DeleteDealCommLog(w http.ResponseWriter, r *http.Request) {
    slog.Info("DeleteDealCommLog called", "method", r.Method, "url", r.URL.Path)
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

    logID, err := strconv.Atoi(chi.URLParam(r, "logId"))
    if err != nil {
        slog.Error("Invalid communication log ID", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid communication log ID")
        return
    }

    existingLog, err := h.commLogService.GetCommLogByID(logID)
    if err != nil {
        slog.Error("Communication log not found", "logID", logID, "error", err)
        respondWithError(w, http.StatusNotFound, "Communication log not found: "+err.Error())
        return
    }

    if existingLog.DealID == nil || *existingLog.DealID != dealID {
        slog.Warn("Communication log does not belong to deal", "logID", logID, "dealID", dealID)
        respondWithError(w, http.StatusNotFound, "Communication log not found for this deal")
        return
    }

    if existingLog.UserID != userID {
        slog.Warn("Unauthorized to delete communication log", "userID", userID, "logUserID", existingLog.UserID)
        respondWithError(w, http.StatusForbidden, "You can only delete your own communication logs")
        return
    }

    if err := h.commLogService.DeleteDealCommLog(logID); err != nil {
        slog.Error("Failed to delete deal communication log", "logID", logID, "error", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to delete deal communication log: "+err.Error())
        return
    }

    slog.Info("Successfully deleted deal communication log", "logID", logID, "dealID", dealID)
    respondWithJSON(w, http.StatusOK, map[string]string{"message": "Deal communication log deleted successfully"})
}