package service

import (
	"context"
	"crm-project/internal/config" // Import config
	"crm-project/internal/dto"
	"crm-project/internal/models"
	"crm-project/internal/repository/postgres"
	"crm-project/internal/util"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   *postgres.UserRepo
	cfg    *config.Config // Add config here
	logger *slog.Logger
}

func NewUserService(repo *postgres.UserRepo, cfg *config.Config, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, cfg: cfg, logger: logger}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (int, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return 0, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can create users via this method.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for CreateUser", "user_id", claims.UserID, "role_id", claims.RoleID)
		return 0, fmt.Errorf("forbidden: only managers can create users")
	}

	if err := util.ValidateStruct(req); err != nil {
		return 0, err
	}

	// Check if username or email already exists
	existingUser, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Error("database error checking for existing user by username", "error", err, "username", req.Username)
		return 0, errors.New("could not verify user existence")
	}
	if existingUser != nil {
		return 0, errors.New("username already taken")
	}
	existingUser, err = s.repo.GetByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Error("database error checking for existing user by email", "error", err, "email", req.Email)
		return 0, errors.New("could not verify user existence")
	}
	if existingUser != nil {
		return 0, errors.New("email already registered")
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
		RoleID:       req.RoleID, // Use the provided RoleID for manager creation
	}

	newUserID, err := s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user in database", "error", err)
		return 0, errors.New("failed to create user")
	}
	s.logger.Info("User created successfully by manager", "manager_id", claims.UserID, "new_user_id", newUserID, "new_user_role", req.RoleID)
	return newUserID, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can view all users.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for GetAllUsers", "user_id", claims.UserID, "role_id", claims.RoleID)
		return nil, fmt.Errorf("forbidden: you do not have permission to view all users")
	}

	return s.repo.GetAll(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return nil, errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can view any user by ID.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for GetUserByID", "user_id", claims.UserID, "role_id", claims.RoleID, "requested_user_id", id)
		return nil, fmt.Errorf("forbidden: you do not have permission to view this user")
	}

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
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can update users.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for UpdateUser", "user_id", claims.UserID, "role_id", claims.RoleID, "target_user_id", id)
		return fmt.Errorf("forbidden: only managers can update users")
	}

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
	s.logger.Info("User updated successfully by manager", "manager_id", claims.UserID, "updated_user_id", id)
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	claims, ok := util.GetClaimsFromContext(ctx)
	if !ok {
		return errors.New("could not retrieve user claims from context")
	}

	// --- PERMISSION CHECK ---
	// Only Reception (Manager) can delete users.
	if claims.RoleID != s.cfg.Roles.ReceptionID {
		s.logger.Warn("Permission denied for DeleteUser", "user_id", claims.UserID, "role_id", claims.RoleID, "target_user_id", id)
		return fmt.Errorf("forbidden: only managers can delete users")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with ID %d not found", id)
		}
		return err
	}
	s.logger.Info("User deleted successfully by manager", "manager_id", claims.UserID, "deleted_user_id", id)
	return nil
}
