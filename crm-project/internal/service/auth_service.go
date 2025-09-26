package service

import (
	"context"
	"crm-project/internal/config" // Import config
	"crm-project/internal/dto"
	"crm-project/internal/models" // Import models for User struct
	"crm-project/internal/repository/postgres"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *postgres.UserRepo
	cfg      *config.Config // Store the entire config
	logger   *slog.Logger
}

func NewAuthService(userRepo *postgres.UserRepo, cfg *config.Config, logger *slog.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
		logger:   logger,
	}
}

func (s *AuthService) LoginUser(ctx context.Context, username, password string) (string, int, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		s.logger.Error("database error finding user by username", "error", err, "username", username)
		return "", 0, err
	}
	if user == nil {
		return "", 0, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", 0, errors.New("invalid credentials")
	}

	s.logger.Info("user authenticated successfully", "user_id", user.ID, "username", user.Username)
	token, err := s.generateJWT(user.ID, user.RoleID, user.Username)
	return token, user.RoleID, err
}

func (s *AuthService) RegisterUser(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		s.logger.Error("database error checking for existing user", "error", err, "username", req.Username)
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return nil, errors.New("failed to process password")
	}

	// Create new user with role from request
	newUser := &models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Email:        req.Email,
		RoleID:       req.RoleID, // Assign role from request
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

		newUserID, err := s.userRepo.Create(ctx, *newUser)
	if err != nil {
		s.logger.Error("failed to create user in database", "error", err)
		return nil, errors.New("failed to register user")
	}

	createdUser, err := s.userRepo.GetByID(ctx, newUserID)
	if err != nil {
		s.logger.Error("failed to fetch created user from database", "error", err)
		return nil, errors.New("failed to retrieve registered user")
	}

	s.logger.Info("new user registered successfully", "user_id", createdUser.ID, "username", createdUser.Username, "role_id", createdUser.RoleID)
	return createdUser, nil
}


func (s *AuthService) generateJWT(userID, roleID int, username string) (string, error) {
	claims := &dto.Claims{
		UserID:   userID,
		RoleID:   roleID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Auth.JWTSecret))
}