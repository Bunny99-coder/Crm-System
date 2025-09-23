package service

import (
	"context"
	"crm-project/internal/dto"   // <-- Import for role constants
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"crm-project/internal/util"  // <-- Import for context helpers
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type ContactService struct {
	repo   *postgres.ContactRepo
	logger *slog.Logger
}

func NewContactService(repo *postgres.ContactRepo, logger *slog.Logger) *ContactService {
	return &ContactService{repo: repo, logger: logger}
}

// CreateContact now automatically assigns the logged-in user as the creator.
func (s *ContactService) CreateContact(ctx context.Context, contact models.Contact) (int, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return 0, errors.New("could not retrieve user claims from context")
	}

	// Set the creator of the contact to the currently logged-in user's ID.
	contact.CreatedBy = &claims.UserID
	
	// You can add more validation here if needed (e.g., check for duplicate phone numbers)
	if contact.FirstName == "" || contact.PrimaryPhone == "" {
		return 0, errors.New("first name and primary phone are required")
	}
    
	return s.repo.Create(ctx, contact)
}

// GetAllContacts now intelligently filters the list based on the user's role.
func (s *ContactService) GetAllContacts(ctx context.Context) ([]models.Contact, error) {
    claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}
    
    // If the user is a Sales Agent, only show contacts they created.
    if claims.RoleID == dto.RoleSalesAgent {
		s.logger.Debug("fetching contacts for single sales agent", "user_id", claims.UserID)
        return s.repo.GetAllForUser(ctx, claims.UserID)
    }
    
    // Otherwise (for Reception/Manager), show all contacts.
	s.logger.Debug("fetching all contacts for manager role", "user_id", claims.UserID)
    return s.repo.GetAll(ctx)
}

// GetContactByID remains the same.
func (s *ContactService) GetContactByID(ctx context.Context, id int) (*models.Contact, error) {
	contact, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	if contact == nil { return nil, fmt.Errorf("contact with ID %d not found", id) }
	return contact, nil
}

// UpdateContact now includes a permission check with logging.
func (s *ContactService) UpdateContact(ctx context.Context, id int, contact models.Contact) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok { 
		s.logger.Error("Could not retrieve user claims for contact update", "contact_id", id)
		return errors.New("could not retrieve user claims") 
	}

	s.logger.Debug("Updating contact", "contact_id", id, "user_id", claims.UserID, "role_id", claims.RoleID)

	// First, get the contact we want to update to check its owner.
	existingContact, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch existing contact for update", "contact_id", id, "error", err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("contact with ID %d not found", id)
		}
		return err
	}
	if existingContact == nil { 
		s.logger.Warn("Existing contact not found", "contact_id", id)
		return fmt.Errorf("contact with ID %d not found", id) 
	}

	s.logger.Debug("Existing contact fetched", "contact_id", id, "created_by", existingContact.CreatedBy)

	// --- PERMISSION CHECK ---
	// A user can update if they are a Receptionist OR if they are the original creator.
	isAllowed := claims.RoleID == dto.RoleReception || (existingContact.CreatedBy != nil && *existingContact.CreatedBy == claims.UserID)
	
	if !isAllowed {
		s.logger.Warn("Permission denied for contact update", "user_id", claims.UserID, "role_id", claims.RoleID, "contact_id", id, "contact_created_by", existingContact.CreatedBy)
		return fmt.Errorf("forbidden: you do not have permission to update this contact")
	}
	
	contact.ID = id
	// We should preserve the original creator
	contact.CreatedBy = existingContact.CreatedBy

	s.logger.Debug("Permission granted, updating contact", "contact_id", id)

	err = s.repo.Update(ctx, contact)
	if err != nil {
		s.logger.Error("Failed to update contact in repo", "contact_id", id, "error", err)
		if err == sql.ErrNoRows { 
			return fmt.Errorf("contact with ID %d not found during update", id) 
		}
		return err
	}
	s.logger.Info("Successfully updated contact", "contact_id", id, "user_id", claims.UserID)
	return nil
}

// DeleteContact now includes the same permission check with logging.
func (s *ContactService) DeleteContact(ctx context.Context, id int) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok { 
		s.logger.Error("Could not retrieve user claims for contact deletion", "contact_id", id)
		return errors.New("could not retrieve user claims") 
	}

	s.logger.Debug("Deleting contact", "contact_id", id, "user_id", claims.UserID, "role_id", claims.RoleID)

	existingContact, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch existing contact for deletion", "contact_id", id, "error", err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("contact with ID %d not found", id)
		}
		return err
	}
	if existingContact == nil { 
		s.logger.Warn("Existing contact not found", "contact_id", id)
		return fmt.Errorf("contact with ID %d not found", id) 
	}

	s.logger.Debug("Existing contact fetched", "contact_id", id, "created_by", existingContact.CreatedBy)

	// --- PERMISSION CHECK ---
	isAllowed := claims.RoleID == dto.RoleReception || (existingContact.CreatedBy != nil && *existingContact.CreatedBy == claims.UserID)

	if !isAllowed {
		s.logger.Warn("Permission denied for contact deletion", "user_id", claims.UserID, "role_id", claims.RoleID, "contact_id", id, "contact_created_by", existingContact.CreatedBy)
		return fmt.Errorf("forbidden: you do not have permission to delete this contact")
	}
	
	s.logger.Debug("Permission granted, deleting contact", "contact_id", id)

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete contact in repo", "contact_id", id, "error", err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("contact with ID %d not found", id)
		}
		// Also check for foreign key constraint errors
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			return fmt.Errorf("cannot delete contact: it is linked to existing leads or deals")
		}
		return err
	}
	s.logger.Info("Successfully deleted contact", "contact_id", id, "user_id", claims.UserID)
	return nil
}