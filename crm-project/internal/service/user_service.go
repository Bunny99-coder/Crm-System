// Replace the contents of internal/service/user_service.go
package service

import (
	"context"
	"crm-project/internal/dto"
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"crm-project/internal/util"
	"database/sql"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   *postgres.UserRepo
	logger *slog.Logger
}

func NewUserService(repo *postgres.UserRepo, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (int, error) {
	if err := util.ValidateStruct(req); err != nil {
		return 0, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password during user creation", "error", err)
		return 0, fmt.Errorf("internal server error")
	}

	user := models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Email:        req.Email,
		RoleID:       req.RoleID,
	}

	return s.repo.Create(ctx, user)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req dto.UpdateUserRequest) error {
	if err := util.ValidateStruct(req); err != nil {
		return err
	}

	// Ensure the user exists before proceeding
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// This handles both db errors and the "not found" case from the repo
		return fmt.Errorf("user with ID %d not found", id)
	}

	user := models.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		RoleID:   req.RoleID,
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error("failed to hash password during user update", "error", err)
			return fmt.Errorf("internal server error")
		}
		user.PasswordHash = string(hashedPassword)
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with ID %d not found during update", id)
		}
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with ID %d not found", id)
		}
		return err
	}
	return nil
}