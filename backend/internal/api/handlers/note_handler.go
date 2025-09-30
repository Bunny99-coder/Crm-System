package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
    "log/slog"

	"time"
	"crm-project/internal/models"
	"crm-project/internal/service"
	"crm-project/internal/util"
		"github.com/go-chi/chi/v5"

)

type NoteHandler struct {
	noteService *service.NoteService
}

func NewNoteHandler(noteService *service.NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

// CreateNoteRequest represents the request body for creating a note
type CreateNoteRequest struct {
	Content   string `json:"content"`
	ContactID *int   `json:"contact_id,omitempty"`
	LeadID    *int   `json:"lead_id,omitempty"`
	DealID    *int   `json:"deal_id,omitempty"`
}

// UpdateNoteRequest represents the request body for updating a note
type UpdateNoteRequest struct {
	Content string `json:"content"`
}

// NoteResponse represents the response structure for notes
type NoteResponse struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	UserID    int    `json:"user_id"`
	ContactID *int   `json:"contact_id,omitempty"`
	LeadID    *int   `json:"lead_id,omitempty"`
	DealID    *int   `json:"deal_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Helper function to get user ID from context
func getUserIDFromContext(ctx context.Context) (int, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok || claims == nil {
		return 0, fmt.Errorf("authentication required")
	}
	return claims.UserID, nil
}

// Helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Helper function to respond with error
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

// Helper function to extract URL parameter
func getURLParam(r *http.Request, param string) string {
	// Implement based on your router (chi, gorilla/mux, etc.)
	// Example for chi: return chi.URLParam(r, param)
	// Example for gorilla/mux: vars := mux.Vars(r); return vars[param]
	
	// For now, using a simple implementation
	// You should replace this with your router's specific method
	return chi.URLParam(r, param)
}

// Helper function to convert models.Note to NoteResponse
func convertNoteToResponse(note *models.Note) NoteResponse {
	response := NoteResponse{
		ID:        note.ID,
		Content:   note.Content,
		UserID:    note.UserID,
		ContactID: note.ContactID,
		LeadID:    note.LeadID,
		DealID:    note.DealID,
	}
	
	// Format timestamps if they exist
	if !note.CreatedAt.IsZero() {
		response.CreatedAt = note.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if !note.UpdatedAt.IsZero() {
		response.UpdatedAt = note.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	
	return response
}

// CreateNote handles POST /api/v1/contacts/{contactId}/notes
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	contactID, err := strconv.Atoi(getURLParam(r, "contactId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	note := &models.Note{
		UserID:    userID,
		ContactID: &contactID,
		Content:   req.Content,
	}

	if err := h.noteService.CreateNote(note); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create note: "+err.Error())
		return
	}

	// Get the created note with complete data including timestamps
	createdNote, err := h.noteService.GetNoteByID(note.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get created note: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, convertNoteToResponse(createdNote))
}

// GetContactNotes handles GET /api/v1/contacts/{contactId}/notes
func (h *NoteHandler) GetContactNotes(w http.ResponseWriter, r *http.Request) {
	contactID, err := strconv.Atoi(getURLParam(r, "contactId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	notes, err := h.noteService.GetNotesByContactID(contactID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get notes: "+err.Error())
		return
	}

	// Convert notes to response format
	noteResponses := make([]NoteResponse, len(notes))
	for i, note := range notes {
		noteResponses[i] = convertNoteToResponse(&note)
	}

	respondWithJSON(w, http.StatusOK, noteResponses)
}

// GetNoteByID handles GET /api/v1/contacts/{contactId}/notes/{noteId}
func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	noteID, err := strconv.Atoi(getURLParam(r, "noteId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	contactID, err := strconv.Atoi(getURLParam(r, "contactId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	note, err := h.noteService.GetNoteByID(noteID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Note not found: "+err.Error())
		return
	}

	// Verify that the note belongs to the specified contact
	if note.ContactID == nil || *note.ContactID != contactID {
		respondWithError(w, http.StatusNotFound, "Note not found for this contact")
		return
	}

	respondWithJSON(w, http.StatusOK, convertNoteToResponse(note))
}

// UpdateNote handles PUT /api/v1/contacts/{contactId}/notes/{noteId}
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	noteID, err := strconv.Atoi(getURLParam(r, "noteId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	contactID, err := strconv.Atoi(getURLParam(r, "contactId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	// Verify the note exists and belongs to the contact
	existingNote, err := h.noteService.GetNoteByID(noteID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Note not found: "+err.Error())
		return
	}

	if existingNote.ContactID == nil || *existingNote.ContactID != contactID {
		respondWithError(w, http.StatusNotFound, "Note not found for this contact")
		return
	}

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	note := &models.Note{
		ID:      noteID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := h.noteService.UpdateNote(note); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update note: "+err.Error())
		return
	}

	// Get the updated note to return complete data
	updatedNote, err := h.noteService.GetNoteByID(noteID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get updated note: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, convertNoteToResponse(updatedNote))
}

// DeleteNote handles DELETE /api/v1/contacts/{contactId}/notes/{noteId}
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	noteID, err := strconv.Atoi(getURLParam(r, "noteId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	contactID, err := strconv.Atoi(getURLParam(r, "contactId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	// Verify the note exists and belongs to the contact
	existingNote, err := h.noteService.GetNoteByID(noteID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Note not found: "+err.Error())
		return
	}

	if existingNote.ContactID == nil || *existingNote.ContactID != contactID {
		respondWithError(w, http.StatusNotFound, "Note not found for this contact")
		return
	}

	// Verify the user owns the note (optional security check)
	if existingNote.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can only delete your own notes")
		return
	}

	if err := h.noteService.DeleteNote(noteID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete note: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Note deleted successfully"})
}

// GetUserNotes handles GET /api/v1/users/{userId}/notes
func (h *NoteHandler) GetUserNotes(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(getURLParam(r, "userId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	notes, err := h.noteService.GetNotesByUserID(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get user notes: "+err.Error())
		return
	}

	// Convert notes to response format
	noteResponses := make([]NoteResponse, len(notes))
	for i, note := range notes {
		noteResponses[i] = convertNoteToResponse(&note)
	}

	respondWithJSON(w, http.StatusOK, noteResponses)
}


















// CreateDealNote handles POST /api/v1/deals/{dealId}/notes
func (h *NoteHandler) CreateDealNote(w http.ResponseWriter, r *http.Request) {
	slog.Info("CreateDealNote called",
		"method", r.Method,
		"url", r.URL.Path,
		"dealID", chi.URLParam(r, "dealId"))

	// Temporarily bypass userID check to test route
	// userID, err := getUserIDFromContext(r.Context())
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Authentication required")
	// 	return
	// }
	// slog.Debug("Authenticated user", "userID", userID)

	// Get dealId from URL
	dealIDStr := chi.URLParam(r, "dealId")
	dealIDInt, err := strconv.Atoi(dealIDStr)
	if err != nil {
		slog.Error("Invalid deal ID", "dealID", dealIDStr, "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}
	dealID := &dealIDInt // Pointer, because models.Note.DealID is *int

	// Parse request body
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			slog.Warn("Failed to close request body", "error", err)
		}
	}()
	slog.Debug("Parsed request body", "content", req.Content)

	// Validate request
	if req.Content == "" {
		slog.Error("Content is empty")
		respondWithError(w, http.StatusBadRequest, "Content is required")
		return
	}

	now := time.Now()
	note := &models.Note{
		Content:   req.Content,
		DealID:    dealID,
		UserID:    0, // Temporary placeholder since auth is bypassed
		CreatedAt: now,
		UpdatedAt: &now,
	}

	// Insert note into database
	if err := h.noteService.CreateNote(note); err != nil {
		slog.Error("Failed to create deal note", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create deal note: "+err.Error())
		return
	}

	// Fetch the created note
	createdNote, err := h.noteService.GetNoteByID(note.ID)
	if err != nil {
		slog.Error("Failed to fetch created note", "noteID", note.ID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch created note: "+err.Error())
		return
	}

	slog.Info("Successfully created deal note", "noteID", createdNote.ID, "dealID", dealIDInt)
	respondWithJSON(w, http.StatusCreated, convertNoteToResponse(createdNote))
}






























// UpdateDealNote handles PUT /deals/{dealId}/notes/{noteId}
func (h *NoteHandler) UpdateDealNote(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	noteID, err := strconv.Atoi(getURLParam(r, "noteId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	dealID, err := strconv.Atoi(getURLParam(r, "dealId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	existingNote, err := h.noteService.GetNoteByID(noteID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Note not found: "+err.Error())
		return
	}

	if existingNote.DealID == nil || *existingNote.DealID != dealID {
		respondWithError(w, http.StatusNotFound, "Note not found for this deal")
		return
	}

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	note := &models.Note{
		ID:      noteID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := h.noteService.UpdateDealNote(note); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update deal note: "+err.Error())
		return
	}

	updatedNote, err := h.noteService.GetNoteByID(noteID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get updated deal note: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, convertNoteToResponse(updatedNote))
}

// DeleteDealNote handles DELETE /deals/{dealId}/notes/{noteId}
func (h *NoteHandler) DeleteDealNote(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	noteID, err := strconv.Atoi(getURLParam(r, "noteId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	dealID, err := strconv.Atoi(getURLParam(r, "dealId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	existingNote, err := h.noteService.GetNoteByID(noteID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Note not found: "+err.Error())
		return
	}

	if existingNote.DealID == nil || *existingNote.DealID != dealID {
		respondWithError(w, http.StatusNotFound, "Note not found for this deal")
		return
	}

	if existingNote.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can only delete your own notes")
		return
	}

	if err := h.noteService.DeleteDealNote(noteID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete deal note: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Deal note deleted successfully"})
}




// GetDealNoteByID handles GET /deals/{dealId}/notes/{noteId}
func (h *NoteHandler) GetDealNoteByID(w http.ResponseWriter, r *http.Request) {
    dealIDStr := chi.URLParam(r, "dealId")
    noteIDStr := chi.URLParam(r, "noteId")

    dealID, err := strconv.Atoi(dealIDStr)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
        return
    }

    noteID, err := strconv.Atoi(noteIDStr)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid note ID")
        return
    }

    note, err := h.noteService.GetNoteByID(noteID)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Note not found: "+err.Error())
        return
    }

    // Verify the note belongs to this deal
    if note.DealID == nil || *note.DealID != dealID {
        respondWithError(w, http.StatusNotFound, "Note not found for this deal")
        return
    }

    respondWithJSON(w, http.StatusOK, convertNoteToResponse(note))
}







// GetDealNotes handles GET /api/v1/deals/{dealId}/notes
func (h *NoteHandler) GetDealNotes(w http.ResponseWriter, r *http.Request) {
	slog.Info("GetDealNotes called", "method", r.Method, "url", r.URL.Path)
	dealID, err := strconv.Atoi(chi.URLParam(r, "dealId"))
	if err != nil {
		slog.Error("Invalid deal ID", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid deal ID")
		return
	}

	notes, err := h.noteService.GetNotesByDealID(dealID)
	if err != nil {
		slog.Error("Failed to get deal notes", "dealID", dealID, "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get deal notes: "+err.Error())
		return
	}

	noteResponses := make([]NoteResponse, len(notes))
	for i, note := range notes {
		noteResponses[i] = convertNoteToResponse(&note)
	}

	slog.Info("Successfully retrieved deal notes", "dealID", dealID, "count", len(notes))
	respondWithJSON(w, http.StatusOK, noteResponses)
}




