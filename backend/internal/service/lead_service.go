package service

import (
	"context"
	"crm-project/internal/config" // Import config
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
	cfg          *config.Config // Add config here
	logger       *slog.Logger
}

func NewLeadService(lr *postgres.LeadRepo, cr *postgres.ContactRepo, ur *postgres.UserRepo, pr *postgres.PropertyRepo, cfg *config.Config, logger *slog.Logger) *LeadService {
	return &LeadService{leadRepo: lr, contactRepo: cr, userRepo: ur, propertyRepo: pr, cfg: cfg, logger: logger}
}

// THIS METHOD IS NOW ROLE-AWARE
func (s *LeadService) GetAllLeads(ctx context.Context) ([]models.Lead, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	// If the user is a Sales Agent, only show leads assigned to them.
	if claims.RoleID == s.cfg.Roles.SalesAgentID { // Use cfg.Roles.SalesAgentID
		s.logger.Debug("fetching leads for single sales agent", "user_id", claims.UserID)
		return s.leadRepo.GetAllLeadsForUser(ctx, claims.UserID)
	}

	// For other roles (like Manager/Reception), show all leads.
	s.logger.Debug("fetching all leads for manager role", "user_id", claims.UserID)
	return s.leadRepo.GetAll(ctx)
}

// THIS METHOD NOW HAS ADVANCED VALIDATION
func (s *LeadService) CreateLead(ctx context.Context, l models.Lead) (int, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return 0, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can create leads.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for CreateLead", "user_id", claims.UserID, "role_id", claims.RoleID)
		return 0, fmt.Errorf("forbidden: only managers can create leads")
	}

	// --- Basic & Foreign Key Validation ---
	if l.ContactID <= 0 || l.SourceID <= 0 || l.StatusID <= 0 || l.AssignedTo <= 0 {
		return 0, errors.New("contact_id, source_id, status_id, and assigned_to are required fields")
	}
	if _, err := s.contactRepo.GetByID(ctx, l.ContactID); err != nil {
		return 0, fmt.Errorf("invalid contact_id: %d", l.ContactID)
	}
	if _, err := s.userRepo.GetByID(ctx, l.AssignedTo); err != nil {
		return 0, fmt.Errorf("invalid assigned_to user_id: %d", l.AssignedTo)
	}

	// --- "One Open Lead per Contact" VALIDATION ---
	hasOpenLead, err := s.leadRepo.CheckForOpenLeadByContactID(ctx, l.ContactID)
	if err != nil {
		s.logger.Error("failed to check for existing open lead", "error", err, "contact_id", l.ContactID)
		return 0, errors.New("could not verify lead status")
	}
	if hasOpenLead {
		return 0, fmt.Errorf("this contact already has an active lead")
	}

	// --- "Property Exclusivity" VALIDATION ---
	if l.PropertyID != nil && *l.PropertyID > 0 {
		if _, err := s.propertyRepo.GetByID(ctx, *l.PropertyID); err != nil {
			return 0, fmt.Errorf("invalid property_id: %d", *l.PropertyID)
		}
		isTaken, err := s.propertyRepo.IsPropertyInOpenLeadOrDeal(ctx, *l.PropertyID)
		if err != nil {
			s.logger.Error("failed to check property availability", "error", err, "property_id", *l.PropertyID)
			return 0, errors.New("could not verify property availability")
		}
		if isTaken {
			return 0, fmt.Errorf("property is already part of an active lead or deal")
		}
	}

	return s.leadRepo.Create(ctx, l)
}


func (s *LeadService) GetLeadByID(ctx context.Context, id int) (*models.Lead, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		s.logger.Error("Could not retrieve user claims for GetLeadByID", "lead_id", id)
		return nil, errors.New("could not retrieve user claims")
	}

	lead, err := s.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if lead == nil {
		return nil, fmt.Errorf("lead with ID %d not found", id)
	}

	// --- PERMISSION CHECK ---
	// A user can view if they are a Receptionist OR if they are the assigned sales agent.
	isAllowed := claims.RoleID == s.cfg.Roles.ReceptionID || (lead.AssignedTo == claims.UserID)

	if !isAllowed {
		s.logger.Warn("Permission denied for GetLeadByID", "user_id", claims.UserID, "role_id", claims.RoleID, "lead_id", id, "lead_assigned_to", lead.AssignedTo)
		return nil, fmt.Errorf("forbidden: you do not have permission to view this lead")
	}

	return lead, nil
}

func (s *LeadService) UpdateLead(ctx context.Context, id int, l models.Lead) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		s.logger.Error("Could not retrieve user claims for UpdateLead", "lead_id", id)
		return errors.New("could not retrieve user claims")
	}

	existingLead, err := s.leadRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch existing lead for update", "lead_id", id, "error", err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("lead with ID %d not found", id)
		}
		return err
	}
	if existingLead == nil {
		s.logger.Warn("Existing lead not found", "lead_id", id)
		return fmt.Errorf("lead with ID %d not found", id)
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can update leads. Sales agents cannot manage leads.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for UpdateLead", "user_id", claims.UserID, "role_id", claims.RoleID, "lead_id", id)
		return fmt.Errorf("forbidden: only managers can update leads")
	}

	l.ID = id
	err = s.leadRepo.Update(ctx, l)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("lead with ID %d not found during update", id)
		}
		return err
	}
	return nil
}

func (s *LeadService) DeleteLead(ctx context.Context, id int) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		s.logger.Error("Could not retrieve user claims for DeleteLead", "lead_id", id)
		return errors.New("could not retrieve user claims")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can delete leads. Sales agents cannot manage leads.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for DeleteLead", "user_id", claims.UserID, "role_id", claims.RoleID, "lead_id", id)
		return fmt.Errorf("forbidden: only managers can delete leads")
	}

	err := s.leadRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("lead with ID %d not found", id)
		}
		return err
	}
	return nil
}
