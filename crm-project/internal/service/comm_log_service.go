package service

import (
    "errors"
    "fmt"

    "crm-project/internal/models"
    "crm-project/internal/repository/postgres"
)

type CommLogService struct {
    commLogRepo postgres.CommLogRepository
}

func NewCommLogService(commLogRepo postgres.CommLogRepository) *CommLogService {
    return &CommLogService{commLogRepo: commLogRepo}
}

// CreateCommLog creates a new communication log
func (s *CommLogService) CreateCommLog(log *models.CommLog) error {
    if log.ContactID != nil && *log.ContactID <= 0 {
        return errors.New("invalid contact ID")
    }
    if log.UserID == 0 {
        return errors.New("user ID is required")
    }
    if log.InteractionDate.IsZero() {
        return errors.New("interaction date is required")
    }
    if log.InteractionType == "" {
        return errors.New("interaction type is required")
    }

    return s.commLogRepo.CreateCommLog(log)
}

// GetCommLogByID retrieves a communication log by ID
func (s *CommLogService) GetCommLogByID(id int) (*models.CommLog, error) {
    if id <= 0 {
        return nil, errors.New("invalid communication log ID")
    }

    return s.commLogRepo.GetCommLogByID(id)
}

// GetCommLogsByDealID retrieves all communication logs for a specific deal
func (s *CommLogService) GetCommLogsByDealID(dealID int) ([]models.CommLog, error) {
    if dealID <= 0 {
        return nil, errors.New("invalid deal ID")
    }

    return s.commLogRepo.GetCommLogsByDealID(dealID)
}

// GetCommLogsByContactID retrieves all communication logs for a specific contact
func (s *CommLogService) GetCommLogsByContactID(contactID int) ([]models.CommLog, error) {
    if contactID <= 0 {
        return nil, errors.New("invalid contact ID")
    }

    return s.commLogRepo.GetCommLogsByContactID(contactID)
}

// GetAllCommLogs retrieves all communication logs
func (s *CommLogService) GetAllCommLogs() ([]models.CommLog, error) {
    return s.commLogRepo.GetAllCommLogs()
}

// UpdateCommLog updates an existing communication log
func (s *CommLogService) UpdateCommLog(log *models.CommLog) error {
    if log.ID <= 0 {
        return errors.New("invalid communication log ID")
    }
    if log.ContactID != nil && *log.ContactID <= 0 {
        return errors.New("invalid contact ID")
    }
    if log.UserID == 0 {
        return errors.New("user ID is required")
    }
    if log.InteractionType == "" {
        return errors.New("interaction type is required")
    }

    existingLog, err := s.commLogRepo.GetCommLogByID(log.ID)
    if err != nil {
        return fmt.Errorf("failed to verify communication log existence: %w", err)
    }

    if existingLog.UserID != log.UserID {
        return errors.New("unauthorized to update this communication log")
    }

    return s.commLogRepo.UpdateCommLog(log)
}

// DeleteCommLog soft deletes a communication log
func (s *CommLogService) DeleteCommLog(id int) error {
    if id <= 0 {
        return errors.New("invalid communication log ID")
    }

    return s.commLogRepo.DeleteCommLog(id)
}

// GetCommLogsForUser retrieves communication logs for a specific user
func (s *CommLogService) GetCommLogsForUser(userID int) ([]models.CommLog, error) {
    if userID <= 0 {
        return nil, errors.New("invalid user ID")
    }

    return s.commLogRepo.GetCommLogsForUser(userID)
}

// CreateDealCommLog creates a new communication log for a deal
func (s *CommLogService) CreateDealCommLog(log *models.CommLog) error {
    if log.ContactID != nil && *log.ContactID <= 0 {
        return errors.New("invalid contact ID")
    }
    if log.UserID == 0 {
        return errors.New("user ID is required")
    }
    if log.DealID == nil || *log.DealID <= 0 {
        return errors.New("deal ID is required")
    }
    if log.InteractionDate.IsZero() {
        return errors.New("interaction date is required")
    }
    if log.InteractionType == "" {
        return errors.New("interaction type is required")
    }

    return s.commLogRepo.CreateCommLog(log)
}

// UpdateDealCommLog updates a communication log for a deal
func (s *CommLogService) UpdateDealCommLog(log *models.CommLog) error {
    if log.ID <= 0 {
        return errors.New("invalid communication log ID")
    }
    if log.ContactID != nil && *log.ContactID <= 0 {
        return errors.New("invalid contact ID")
    }
    if log.InteractionType == "" {
        return errors.New("interaction type is required")
    }

    existingLog, err := s.commLogRepo.GetCommLogByID(log.ID)
    if err != nil {
        return fmt.Errorf("failed to verify communication log existence: %w", err)
    }

    if existingLog.DealID == nil || *existingLog.DealID != *log.DealID {
        return errors.New("communication log not found for this deal")
    }

    return s.commLogRepo.UpdateCommLog(log)
}

// DeleteDealCommLog deletes a communication log for a deal
func (s *CommLogService) DeleteDealCommLog(id int) error {
    if id <= 0 {
        return errors.New("invalid communication log ID")
    }

    return s.commLogRepo.DeleteCommLog(id)
}

// CreateContactCommLog creates a new communication log for a contact
func (s *CommLogService) CreateContactCommLog(log *models.CommLog) error {
    if log.ContactID == nil || *log.ContactID <= 0 {
        return errors.New("contact ID is required")
    }
    if log.UserID == 0 {
        return errors.New("user ID is required")
    }
    if log.InteractionDate.IsZero() {
        return errors.New("interaction date is required")
    }
    if log.InteractionType == "" {
        return errors.New("interaction type is required")
    }

    return s.commLogRepo.CreateCommLog(log)
}

// UpdateContactCommLog updates a communication log for a contact
func (s *CommLogService) UpdateContactCommLog(log *models.CommLog) error {
    if log.ID <= 0 {
        return errors.New("invalid communication log ID")
    }
    if log.ContactID == nil || *log.ContactID <= 0 {
        return errors.New("contact ID is required")
    }
    if log.InteractionType == "" {
        return errors.New("interaction type is required")
    }

    existingLog, err := s.commLogRepo.GetCommLogByID(log.ID)
    if err != nil {
        return fmt.Errorf("failed to verify communication log existence: %w", err)
    }

    if existingLog.ContactID == nil || log.ContactID == nil || *existingLog.ContactID != *log.ContactID {
        return errors.New("communication log not found for this contact")
    }

    return s.commLogRepo.UpdateCommLog(log)
}

// DeleteContactCommLog deletes a communication log for a contact
func (s *CommLogService) DeleteContactCommLog(id int) error {
    if id <= 0 {
        return errors.New("invalid communication log ID")
    }

    return s.commLogRepo.DeleteCommLog(id)
}