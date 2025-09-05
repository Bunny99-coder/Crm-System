// Replace the contents of internal/service/property_service.go
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

type PropertyService struct {
	repo   *postgres.PropertyRepo
	logger *slog.Logger
}

func NewPropertyService(repo *postgres.PropertyRepo, logger *slog.Logger) *PropertyService {
	return &PropertyService{repo: repo, logger: logger}
}

func (s *PropertyService) CreateProperty(ctx context.Context, p models.Property) (int, error) {
	if p.Name == "" || p.Price <= 0 || p.SiteID <= 0 || p.PropertyTypeID <= 0 {
		return 0, errors.New("name, price, site_id, and property_type_id are required fields")
	}
	// In a real app, you would also validate that site_id and property_type_id exist.
	return s.repo.Create(ctx, p)
}

func (s *PropertyService) GetAllProperties(ctx context.Context) ([]models.Property, error) {
	return s.repo.GetAll(ctx)
}

func (s *PropertyService) GetPropertyByID(ctx context.Context, id int) (*models.Property, error) {
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
	_, err := s.GetPropertyByID(ctx, id)
	if err != nil {
		return err
	}
	p.ID = id
	if p.Name == "" || p.Price <= 0 || p.SiteID <= 0 || p.PropertyTypeID <= 0 {
		return errors.New("name, price, site_id, and property_type_id are required fields")
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
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("property with ID %d not found", id)
		}
		return err
	}
	return nil
}