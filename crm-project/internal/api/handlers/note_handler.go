// Replace the contents of internal/api/handlers/note_handler.go
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

type NoteHandler struct {
	service *service.NoteService
	logger  *slog.Logger
}

func NewNoteHandler(s *service.NoteService, logger *slog.Logger) *NoteHandler {
	return &NoteHandler{service: s, logger: logger}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newNote models.Note
	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		h.logger.Warn("invalid create note request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newID, err := h.service.CreateNote(ctx, newNote)
	if err != nil {
		h.logger.Error("failed to create note", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.logger.Info("note created successfully", "note_id", newID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": newID})
}



// in internal/api/handlers/note_handler.go
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	notes, err := h.service.GetAllNotes(ctx)
	if err != nil {
		h.logger.Error("failed to get all notes", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.logger.Debug("retrieved all notes", "count", len(notes))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}



func (h *NoteHandler) GetNotesByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	notes, err := h.service.GetNotesByUser(ctx, userID)
	if err != nil {
		h.logger.Error("failed to get notes for user", "user_id", userID, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Debug("retrieved all notes for user", "user_id", userID, "count", len(notes))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "noteId"))
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}
	note, err := h.service.GetNoteByID(ctx, id)
	if err != nil {
		h.logger.Warn("note not found", "note_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Debug("retrieved note by id", "note_id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "noteId"))
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}
	var n models.Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		h.logger.Warn("invalid update note request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = h.service.UpdateNote(ctx, id, n)
	if err != nil {
		h.logger.Error("failed to update note", "note_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("note updated successfully", "note_id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "noteId"))
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteNote(ctx, id)
	if err != nil {
		h.logger.Error("failed to delete note", "note_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("note deleted successfully", "note_id", id)
	w.WriteHeader(http.StatusNoContent)
}