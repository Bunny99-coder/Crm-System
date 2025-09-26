package service

import (
	"context"
	"crm-project/internal/config" // Import config
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"crm-project/internal/util" // Import util for claims
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type DealService struct {
	dealRepo     *postgres.DealRepo
	leadRepo     *postgres.LeadRepo
	propertyRepo *postgres.PropertyRepo
	cfg          *config.Config // Add config here
	logger       *slog.Logger
}

func NewDealService(dr *postgres.DealRepo, lr *postgres.LeadRepo, pr *postgres.PropertyRepo, cfg *config.Config, logger *slog.Logger) *DealService {
	return &DealService{dealRepo: dr, leadRepo: lr, propertyRepo: pr, cfg: cfg, logger: logger}
}

// THIS METHOD NOW HAS ADVANCED VALIDATION AND ROLE-AWARENESS
func (s *DealService) CreateDeal(ctx context.Context, d models.Deal) (int, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return 0, errors.New("could not retrieve user claims from context")
	}

	// --- Deal Integrity Validation ---
	lead, err := s.leadRepo.GetByID(ctx, d.LeadID)
	if err != nil || lead == nil {
		return 0, fmt.Errorf("invalid lead_id: %d", d.LeadID)
	}
	if lead.PropertyID == nil {
		return 0, errors.New("cannot create a deal from a lead that is not linked to a property")
	}
	if d.PropertyID != *lead.PropertyID {
		return 0, fmt.Errorf("deal property ID (%d) does not match the lead's property ID (%d)", d.PropertyID, *lead.PropertyID)
	}
	if d.DealAmount <= 0 {
		return 0, errors.New("deal amount must be positive")
	}

	// Set the creator of the deal to the currently logged-in user's ID.
	d.CreatedBy = claims.UserID

	newID, err := s.dealRepo.Create(ctx, d)
	if err != nil {
		return 0, err
	}

	// --- Automatic Property Status Update ---
	if d.DealStatus == "Closed-Won" {
		if err := s.updatePropertyStatusOnDealClose(ctx, d.PropertyID); err != nil {
			s.logger.Warn("deal created, but failed to update property status", "deal_id", newID, "error", err)
		}
	}
	return newID, nil
}

// THIS METHOD NOW HAS ADVANCED LOGIC AND ROLE-AWARENESS
func (s *DealService) UpdateDeal(ctx context.Context, id int, d models.Deal) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		s.logger.Error("Could not retrieve user claims for deal update", "deal_id", id)
		return errors.New("could not retrieve user claims")
	}

	existingDeal, err := s.dealRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch existing deal for update", "deal_id", id, "error", err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("deal with ID %d not found", id)
		}
		return err
	}
	if existingDeal == nil {
		s.logger.Warn("Existing deal not found", "deal_id", id)
		return fmt.Errorf("deal with ID %d not found", id)
	}

	// --- PERMISSION CHECK ---
	// A user can update if they are a Receptionist OR if they are the original creator.
	isAllowed := claims.RoleID == s.cfg.Roles.ReceptionID || (existingDeal.CreatedBy == claims.UserID)

	if !isAllowed {
		s.logger.Warn("Permission denied for deal update", "user_id", claims.UserID, "role_id", claims.RoleID, "deal_id", id, "deal_created_by", existingDeal.CreatedBy)
		return fmt.Errorf("forbidden: you do not have permission to update this deal")
	}

	d.ID = id
	// Preserve original creator
	d.CreatedBy = existingDeal.CreatedBy

	err = s.dealRepo.Update(ctx, d)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("deal with ID %d not found during update", id)
		}
		return err
	}

	// --- Automatic Property Status Update (only if status changes to Closed-Won) ---
	if d.DealStatus == "Closed-Won" && existingDeal.DealStatus != "Closed-Won" {
		if err := s.updatePropertyStatusOnDealClose(ctx, d.PropertyID); err != nil {
			s.logger.Warn("deal updated, but failed to update property status", "deal_id", id, "error", err)
		}
	}
	return nil
}

// New helper function
func (s *DealService) updatePropertyStatusOnDealClose(ctx context.Context, propertyID int) error {
	s.logger.Info("deal closed, attempting to update property status to Sold", "property_id", propertyID)
	property, err := s.propertyRepo.GetByID(ctx, propertyID)
	if err != nil {
		return fmt.Errorf("could not find property to update: %w", err)
	}

	property.Status = "Sold"
	return s.propertyRepo.Update(ctx, *property)
}

// GetAllDeals now intelligently filters the list based on the user's role.
func (s *DealService) GetAllDeals(ctx context.Context) ([]models.Deal, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	// If the user is a Sales Agent, only show deals they created.
	if claims.RoleID == s.cfg.Roles.SalesAgentID {
		s.logger.Debug("fetching deals for single sales agent", "user_id", claims.UserID)
		return s.dealRepo.GetAllForUser(ctx, claims.UserID)
	}

	// For other roles (like Manager/Reception), show all deals.
	s.logger.Debug("fetching all deals for manager role", "user_id", claims.UserID)
	return s.dealRepo.GetAll(ctx)
}

// GetDealByID now includes a permission check.
func (s *DealService) GetDealByID(ctx context.Context, id int) (*models.Deal, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		s.logger.Error("Could not retrieve user claims for GetDealByID", "deal_id", id)
		return nil, errors.New("could not retrieve user claims")
	}

	deal, err := s.dealRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if deal == nil {
		return nil, fmt.Errorf("deal with ID %d not found", id)
	}

	// --- PERMISSION CHECK ---
	// A user can view if they are a Receptionist OR if they are the original creator.
	isAllowed := claims.RoleID == s.cfg.Roles.ReceptionID || (deal.CreatedBy == claims.UserID)

	if !isAllowed {
		s.logger.Warn("Permission denied for GetDealByID", "user_id", claims.UserID, "role_id", claims.RoleID, "deal_id", id, "deal_created_by", deal.CreatedBy)
		return nil, fmt.Errorf("forbidden: you do not have permission to view this deal")
	}

	return deal, nil
}

// DeleteDeal now includes a permission check.
func (s *DealService) DeleteDeal(ctx context.Context, id int) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		s.logger.Error("Could not retrieve user claims for DeleteDeal", "deal_id", id)
		return errors.New("could not retrieve user claims")
	}

	existingDeal, err := s.dealRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch existing deal for deletion", "deal_id", id, "error", err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("deal with ID %d not found", id)
		}
		return err
	}
	if existingDeal == nil {
		s.logger.Warn("Existing deal not found", "deal_id", id)
		return fmt.Errorf("deal with ID %d not found", id)
	}

	// --- PERMISSION CHECK ---
	// A user can delete if they are a Receptionist OR if they are the original creator.
	isAllowed := claims.RoleID == s.cfg.Roles.ReceptionID || (existingDeal.CreatedBy == claims.UserID)

	if !isAllowed {
		s.logger.Warn("Permission denied for deal deletion", "user_id", claims.UserID, "role_id", claims.RoleID, "deal_id", id, "deal_created_by", existingDeal.CreatedBy)
		return fmt.Errorf("forbidden: you do not have permission to delete this deal")	}

	err = s.dealRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("deal with ID %d not found", id)
		}
		return err
	}
	return nil
}
