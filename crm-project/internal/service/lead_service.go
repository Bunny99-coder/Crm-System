// Replace the entire contents of your internal/service/lead_service.go file with this.
package service

import (
	"context"
	"crm-project/internal/dto"
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"crm-project/internal/util"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type LeadService struct {
	leadRepo     *postgres.LeadRepo
	contactRepo  *postgres.ContactRepo
	userRepo     *postgres.UserRepo
	propertyRepo *postgres.PropertyRepo
	logger       *slog.Logger
}

func NewLeadService(lr *postgres.LeadRepo, cr *postgres.ContactRepo, ur *postgres.UserRepo, pr *postgres.PropertyRepo, logger *slog.Logger) *LeadService {
	return &LeadService{
		leadRepo:     lr,
		contactRepo:  cr,
		userRepo:     ur,
		propertyRepo: pr,
		logger:       logger,
	}
}

// Renamed to match the handler's expectation
// In internal/service/lead_service.go

func (s *LeadService) CreateLead(ctx context.Context, l models.Lead) (int, error) {
	// --- Basic Validation ---
	if l.ContactID <= 0 || l.SourceID <= 0 || l.StatusID <= 0 || l.AssignedTo <= 0 {
		return 0, errors.New("contact_id, source_id, status_id, and assigned_to are required fields")
	}

	// --- Foreign Key Validation ---
	if _, err := s.contactRepo.GetByID(ctx, l.ContactID); err != nil {
		// This handles both errors and the case where the contact is not found.
		return 0, fmt.Errorf("invalid contact_id: %d", l.ContactID)
	}
	if _, err := s.userRepo.GetByID(ctx, l.AssignedTo); err != nil {
		return 0, fmt.Errorf("invalid assigned_to user_id: %d", l.AssignedTo)
	}
	if l.PropertyID != nil && *l.PropertyID > 0 {
		if _, err := s.propertyRepo.GetByID(ctx, *l.PropertyID); err != nil {
			return 0, fmt.Errorf("invalid property_id: %d", *l.PropertyID)
		}
	}
	
	// --- NEW "One Open Lead per Contact" BUSINESS RULE VALIDATION ---
	hasOpenLead, err := s.leadRepo.CheckForOpenLeadByContactID(ctx, l.ContactID)
	if err != nil {
		// Log the internal database error but return a generic message to the user.
		s.logger.Error("failed to check for existing open lead", "error", err, "contact_id", l.ContactID)
		return 0, errors.New("could not verify lead status, please try again")
	}
	if hasOpenLead {
		// This is a business rule violation, so we return a clear error.
		return 0, fmt.Errorf("this contact already has an active lead and cannot be assigned a new one")
	}
	// --- END OF NEW VALIDATION ---





	if l.PropertyID != nil && *l.PropertyID > 0 {
		isTaken, err := s.propertyRepo.IsPropertyInOpenLeadOrDeal(ctx, *l.PropertyID)
		if err != nil {
			s.logger.Error("failed to check property availability", "error", err, "property_id", *l.PropertyID)
			return 0, errors.New("could not verify property availability")
		}
		if isTaken {
			return 0, fmt.Errorf("property is already part of an active lead or deal")
		}
	}







	
	// If all validation passes, proceed to create the lead.
	return s.leadRepo.Create(ctx, l)
}

// Renamed to match the handler's expectation
func (s *LeadService) GetAllLeads(ctx context.Context) ([]models.Lead, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}
	if claims.RoleID == dto.RoleSalesAgent {
		s.logger.Debug("fetching leads for single sales agent", "user_id", claims.UserID)
		return s.leadRepo.GetAllForUser(ctx, claims.UserID)
	}
	s.logger.Debug("fetching all leads for non-sales-agent role", "user_id", claims.UserID)
	return s.leadRepo.GetAll(ctx)
}

// Renamed to match the handler's expectation
func (s *LeadService) GetLeadByID(ctx context.Context, id int) (*models.Lead, error) {
	lead, err := s.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if lead == nil {
		return nil, fmt.Errorf("lead with ID %d not found", id)
	}
	return lead, nil
}

// Renamed to match the handler's expectation
func (s *LeadService) UpdateLead(ctx context.Context, id int, l models.Lead) error {
	_, err := s.GetLeadByID(ctx, id)
	if err != nil {
		return err
	}
	l.ID = id
	if _, err := s.contactRepo.GetByID(ctx, l.ContactID); err != nil {
		return fmt.Errorf("invalid contact_id on update: %d", l.ContactID)
	}
	if _, err := s.userRepo.GetByID(ctx, l.AssignedTo); err != nil {
		return fmt.Errorf("invalid assigned_to user_id on update: %d", l.AssignedTo)
	}
	err = s.leadRepo.Update(ctx, l)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("lead with ID %d not found during update", id)
		}
		return err
	}
	return nil
}

// Renamed to match the handler's expectation
func (s *LeadService) DeleteLead(ctx context.Context, id int) error {
	err := s.leadRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("lead with ID %d not found", id)
		}
		return err
	}
	return nil
}


