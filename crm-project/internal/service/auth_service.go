// Replace the entire contents of internal/service/auth_service.go with this.
package service

import (
	"context"
	"crm-project/internal/repository/postgres"
	"errors"
	"log/slog"
	"time"
 "crm-project/internal/util" // <-- Add this

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"crm-project/internal/dto" // <-- IMPORTANT: Import the 'api' package for the Claims struct

)

type AuthService struct {
	userRepo  *postgres.UserRepo
	jwtSecret string
	logger    *slog.Logger
}

func NewAuthService(userRepo *postgres.UserRepo, jwtSecret string, logger *slog.Logger) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *AuthService) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		s.logger.Error("database error finding user by username", "error", err, "username", username)
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	s.logger.Info("user authenticated successfully", "user_id", user.ID, "username", user.Username)
	return s.generateJWT(user.ID, user.RoleID)
}

// in internal/service/auth_service.go

func (s *AuthService) generateJWT(userID, roleID int) (string, error) {
	// Use our new, strongly-typed Claims struct from the dto package.
	claims := &dto.Claims{ // <-- This line will now work
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// --- TEMPORARY CHANGE FOR DEBUGGING ---
	hardcodedSecret := "my_super_secret_debug_key_12345"
	s.logger.Debug("Signing token with hardcoded secret", "secret", hardcodedSecret)
	return token.SignedString([]byte(util.HardcodedJWTSecret))
}