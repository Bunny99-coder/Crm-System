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

type PropertyService struct {
	repo   *postgres.PropertyRepo
	cfg    *config.Config // Add config here
	logger *slog.Logger
}

func NewPropertyService(repo *postgres.PropertyRepo, cfg *config.Config, logger *slog.Logger) *PropertyService {
	return &PropertyService{repo: repo, cfg: cfg, logger: logger}
}

func (s *PropertyService) CreateProperty(ctx context.Context, p models.Property) (int, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return 0, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can create properties.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for CreateProperty", "user_id", claims.UserID, "role_id", claims.RoleID)
		return 0, fmt.Errorf("forbidden: only managers can create properties")
	}

	if p.Name == "" || p.Price <= 0 || p.SiteID <= 0 || p.PropertyTypeID <= 0 {
		return 0, errors.New("name, price, site_id, and property_type_id are required fields")
	}

	// Validate that site_id exists.
	siteExists, err := s.repo.SiteExists(ctx, p.SiteID)
	if err != nil {
		return 0, fmt.Errorf("error validating site: %w", err)
	}
	if !siteExists {
		return 0, fmt.Errorf("site with ID %d does not exist", p.SiteID)
	}

	// Validate that property_type_id exists.
	propertyTypeExists, err := s.repo.PropertyTypeExists(ctx, p.PropertyTypeID)
	if err != nil {
		return 0, fmt.Errorf("error validating property type: %w", err)
	}
	if !propertyTypeExists {
		return 0, fmt.Errorf("property type with ID %d does not exist", p.PropertyTypeID)
	}

	return s.repo.Create(ctx, p)
}

func (s *PropertyService) GetAllProperties(ctx context.Context) ([]models.Property, error) {
	// Both roles can view all properties, no specific filtering needed here.
	return s.repo.GetAll(ctx)
}

func (s *PropertyService) GetPropertyByID(ctx context.Context, id int) (*models.Property, error) {
	// Both roles can view properties by ID, no specific filtering needed here.
	property, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, fmt.Errorf("property with ID %d not found", id)
	}
	return property, nil
}

func (s *PropertyService) UpdateProperty(ctx context.Context, id int, p models.Property) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can update properties.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for UpdateProperty", "user_id", claims.UserID, "role_id", claims.RoleID, "property_id", id)
		return fmt.Errorf("forbidden: only managers can update properties")
	}

	_, err := s.GetPropertyByID(ctx, id)
	if err != nil {
		return err
	}
	p.ID = id
	if p.Name == "" || p.Price <= 0 || p.SiteID <= 0 || p.PropertyTypeID <= 0 {
		return errors.New("name, price, site_id, and property_type_id are required fields")
	}

	// Validate that site_id exists.
	siteExists, err := s.repo.SiteExists(ctx, p.SiteID)
	if err != nil {
		return fmt.Errorf("error validating site: %w", err)
	}
	if !siteExists {
		return fmt.Errorf("site with ID %d does not exist", p.SiteID)
	}

	// Validate that property_type_id exists.
	propertyTypeExists, err := s.repo.PropertyTypeExists(ctx, p.PropertyTypeID)
	if err != nil {
		return fmt.Errorf("error validating property type: %w", err)
	}
	if !propertyTypeExists {
		return fmt.Errorf("property type with ID %d does not exist", p.PropertyTypeID)
	}

	err = s.repo.Update(ctx, p)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("property with ID %d not found during update", id)
		}
		return err
	}
	return nil
}

func (s *PropertyService) DeleteProperty(ctx context.Context, id int) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can delete properties.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for DeleteProperty", "user_id", claims.UserID, "role_id", claims.RoleID, "property_id", id)
		return fmt.Errorf("forbidden: only managers can delete properties")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("property with ID %d not found", id)
		}
		return err
	}
	return nil
}
