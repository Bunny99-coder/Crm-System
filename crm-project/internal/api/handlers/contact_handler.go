// File: internal/api/handlers/contact_handler.go
package handlers

import (
	"crm-project/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"crm-project/internal/models"
	"github.com/go-chi/chi/v5"
	"strconv"
	"log/slog"


)

// ContactHandler handles HTTP requests for contacts.

type ContactHandler struct {
	service *service.ContactService
	logger  *slog.Logger // <-- Add logger field
}

// NewContactHandler creates a new ContactHandler.
func NewContactHandler(s *service.ContactService, logger *slog.Logger) *ContactHandler {
	return &ContactHandler{service: s, logger: logger}
}



// GetAllContacts is the handler function for GET /contacts
func (h *ContactHandler) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	// The handler calls the service layer to get the data.
		ctx := r.Context()

	contacts, err := h.service.GetAllContacts(ctx)
	if err != nil {
		h.logger.Error("error getting all contacts", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// The handler's job is to format the response correctly.
	h.logger.Debug("retrieved all contacts", "count", len(contacts))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)

	


}

	// Add this method to your contact_handler.go file

// CreateContact is the handler for POST /contacts
func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

	var newContact models.Contact

	// 1. Decode the incoming JSON from the request body into our Contact struct.
	err := json.NewDecoder(r.Body).Decode(&newContact)
	if err != nil {
		// If there's an error in the JSON, it's a client error.
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Call the service layer with the decoded data.
	newID, err := h.service.CreateContact(ctx, newContact)
	if err != nil {
		// Our service layer validation might return an error.
		// A real app would check the type of error and return a more specific status code.
		log.Printf("Error creating contact: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest) // Send back the validation error message
		return
	}

	// 3. Respond with the ID of the newly created contact.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created is the correct status code for a successful creation.
	
	// We can send a simple JSON response back with the new ID.
	json.NewEncoder(w).Encode(map[string]int{"id": newID})
}


// You will need to add "strconv" to your imports for this file.
// Replace your GetContactByID function with this one.

func (h *ContactHandler) GetContactByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// 1. Get the URL parameter using the CORRECT name "contactId" from our router.
	idStr := chi.URLParam(r, "contactId")

	h.logger.Info("GetContactByID called with idStr:", "id", idStr) // Added log

	// 2. Safely convert the string to an integer and HANDLE THE ERROR.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// This will catch cases where the ID is not a number, or is missing.
		h.logger.Warn("invalid contact ID in URL", "raw_id", idStr, "error", err)
		http.Error(w, "Invalid contact ID format", http.StatusBadRequest)
		return
	}

	// 3. Call the service with the now-validated ID.
	contact, err := h.service.GetContactByID(ctx, id)
	if err != nil {
		h.logger.Warn("contact not found", "contact_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	h.logger.Debug("retrieved contact by id", "contact_id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contact)
}


// UpdateContact is the handler for PUT /contacts/{id}
// Replace the UpdateContact function in contact_handler.go with this:
func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "contactId") // <-- THE FIX
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid contact ID format", http.StatusBadRequest)
		return
	}

	var contactToUpdate models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contactToUpdate); err != nil {
		h.logger.Warn("invalid update contact request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	err = h.service.UpdateContact(ctx, id, contactToUpdate)
	if err != nil {
		h.logger.Error("failed to update contact", "contact_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // Or be more specific based on the error
		return
	}
	
	h.logger.Info("contact updated successfully", "contact_id", id)
	w.WriteHeader(http.StatusNoContent)
}


// Replace the DeleteContact function in contact_handler.go with this:
func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "contactId") // <-- THE FIX
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid contact ID format", http.StatusBadRequest)
		return
	}
	
	err = h.service.DeleteContact(ctx, id)
	if err != nil {
		h.logger.Error("failed to delete contact", "contact_id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	h.logger.Info("contact deleted successfully", "contact_id", id)
	w.WriteHeader(http.StatusNoContent)
}