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
	s.logger.Debug("Attempting to create deal", "incoming_deal", d)

	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		s.logger.Error("Could not retrieve user claims from context for CreateDeal")
		return 0, errors.New("could not retrieve user claims from context")
	}
	s.logger.Debug("Claims retrieved from context", "user_id", claims.UserID, "role_id", claims.RoleID)

	// --- PERMISSION CHECK ---
	// Only Reception can create deals.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for CreateDeal", "user_id", claims.UserID, "role_id", claims.RoleID)
		return 0, fmt.Errorf("forbidden: only receptionists can create deals")
	}

	// --- Deal Integrity Validation ---
	lead, err := s.leadRepo.GetByID(ctx, d.LeadID)
	if err != nil {
		s.logger.Warn("Failed to retrieve lead for deal creation", "lead_id", d.LeadID, "error", err)
		return 0, fmt.Errorf("invalid lead_id: %d, error: %w", d.LeadID, err)
	}
	if lead == nil {
		s.logger.Warn("Lead not found for deal creation", "lead_id", d.LeadID)
		return 0, fmt.Errorf("lead with ID %d not found", d.LeadID)
	}
	s.logger.Debug("Lead retrieved", "lead", lead)

	if lead.PropertyID == nil {
		s.logger.Warn("Lead not linked to a property", "lead_id", d.LeadID)
		return 0, errors.New("cannot create a deal from a lead that is not linked to a property")
	}
	if d.PropertyID != *lead.PropertyID {
		s.logger.Warn("Deal property ID mismatch with lead's property ID", "deal_property_id", d.PropertyID, "lead_property_id", *lead.PropertyID)
		return 0, fmt.Errorf("deal property ID (%d) does not match the lead's property ID (%d)", d.PropertyID, *lead.PropertyID)
	}
	if d.DealAmount <= 0 {
		s.logger.Warn("Deal amount is not positive", "deal_amount", d.DealAmount)
		return 0, errors.New("deal amount must be positive")
	}

	// Set the creator of the deal to the currently logged-in user's ID.
	d.CreatedBy = sql.NullInt64{Int64: int64(claims.UserID), Valid: true}
	s.logger.Debug("CreatedBy set for deal", "created_by", d.CreatedBy.Int64)

	newID, err := s.dealRepo.Create(ctx, d)
	if err != nil {
		s.logger.Error("Failed to create deal in repository", "deal", d, "error", err)
		return 0, fmt.Errorf("failed to create deal: %w", err)
	}
	s.logger.Info("Deal successfully created in repository", "deal_id", newID)

	// --- Automatic Property Status Update ---
	if d.DealStatus == "Closed-Won" {
		if err := s.updatePropertyStatusOnDealClose(ctx, d.PropertyID); err != nil {
			s.logger.Warn("Deal created, but failed to update property status", "deal_id", newID, "error", err)
		}
	}
	return newID, nil
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
func (s *DealService) GetAllDeals(ctx context.Context, userID int, roleID int) ([]models.Deal, error) {
	// If the user is a Sales Agent, only show deals they created.
	if roleID == s.cfg.Roles.SalesAgentID {
		s.logger.Debug("fetching deals for single sales agent", "user_id", userID)
		return s.dealRepo.GetAllForUser(ctx, userID)
	}

	// For other roles (like Manager/Reception), show all deals.
	s.logger.Debug("fetching all deals for manager role", "user_id", userID)
	return s.dealRepo.GetAll(ctx)
}

func (s *DealService) GetDealByID(ctx context.Context, dealID int, userID int, roleID int) (*models.Deal, error) {
	deal, err := s.dealRepo.GetByID(ctx, dealID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("deal with ID %d not found", dealID)
		}
		return nil, fmt.Errorf("failed to get deal by ID: %w", err)
	}
	if deal == nil {
		return nil, fmt.Errorf("deal with ID %d not found", dealID)
	}

	// Permission check: Receptionist can view all, SalesAgent can view their own
	isAllowed := roleID == s.cfg.Roles.ReceptionID || (deal.CreatedBy.Valid && deal.CreatedBy.Int64 == int64(userID))

	if !isAllowed {
		s.logger.Warn("Permission denied for deal viewing", "user_id", userID, "role_id", roleID, "deal_id", dealID, "deal_created_by", deal.CreatedBy)
		return nil, fmt.Errorf("unauthorized to view this deal")
	}

	return deal, nil
}


// UpdateDeal now includes a permission check.
func (s *DealService) UpdateDeal(ctx context.Context, id int, d models.Deal, userID int, roleID int) error {
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
	// Only Reception can update deals.
	if roleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for deal update", "user_id", userID, "role_id", roleID, "deal_id", id, "deal_created_by", existingDeal.CreatedBy)
		return fmt.Errorf("unauthorized")
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

// DeleteDeal now includes a permission check.
func (s *DealService) DeleteDeal(ctx context.Context, id int, userID int, roleID int) error {
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
	// Only Reception can delete deals.
	if roleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for deal deletion", "user_id", userID, "role_id", roleID, "deal_id", id, "deal_created_by", existingDeal.CreatedBy)
		return fmt.Errorf("unauthorized")
	}

	err = s.dealRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("deal with ID %d not found", id)
		}
		return err
	}
	return nil
}