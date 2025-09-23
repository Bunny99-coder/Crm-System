// Replace the entire contents of your internal/service/deal_service.go file.
package service

import (
	"context"
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type DealService struct {
	dealRepo     *postgres.DealRepo
	leadRepo     *postgres.LeadRepo
	propertyRepo *postgres.PropertyRepo
	logger       *slog.Logger
}

func NewDealService(dr *postgres.DealRepo, lr *postgres.LeadRepo, pr *postgres.PropertyRepo, logger *slog.Logger) *DealService {
	return &DealService{dealRepo: dr, leadRepo: lr, propertyRepo: pr, logger: logger}
}

// THIS METHOD NOW HAS ADVANCED VALIDATION
func (s *DealService) CreateDeal(ctx context.Context, d models.Deal) (int, error) {
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

// THIS METHOD NOW HAS ADVANCED LOGIC
func (s *DealService) UpdateDeal(ctx context.Context, id int, d models.Deal) error {
	existingDeal, err := s.GetDealByID(ctx, id)
	if err != nil {
		return err
	}
	d.ID = id
	
	err = s.dealRepo.Update(ctx, d)
	if err != nil {
		if err == sql.ErrNoRows { return fmt.Errorf("deal with ID %d not found", id) }
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

func (s *DealService) GetAllDeals(ctx context.Context) ([]models.Deal, error) {
	return s.dealRepo.GetAll(ctx)
}
func (s *DealService) GetDealByID(ctx context.Context, id int) (*models.Deal, error) {
	deal, err := s.dealRepo.GetByID(ctx, id)
	if err != nil { return nil, err }
	if deal == nil { return nil, fmt.Errorf("deal with ID %d not found", id) }
	return deal, nil
}
func (s *DealService) DeleteDeal(ctx context.Context, id int) error {
	err := s.dealRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows { return fmt.Errorf("deal with ID %d not found", id) }
		return err
	}
	return nil
}