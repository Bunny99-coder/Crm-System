// Replace the entire contents of internal/service/auth_service.go
package service

import (
	"context"
	"crm-project/internal/dto"
	"crm-project/internal/repository/postgres"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	// --- THE FIX IS HERE ---
	// We now pass the user's username into the token generator.
	return s.generateJWT(user.ID, user.RoleID, user.Username)
}

// --- THE FIX IS HERE ---
// The function now accepts the username as an argument.
func (s *AuthService) generateJWT(userID, roleID int, username string) (string, error) {
	claims := &dto.Claims{
		UserID:   userID,
		RoleID:   roleID,
		Username: username, // This line will now work correctly.
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}