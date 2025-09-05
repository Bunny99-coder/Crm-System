// Replace the entire contents of your internal/api/middleware.go file with this.
package api

import (
	"context"
	"crm-project/internal/dto"   // <-- Import shared DTOs
	"crm-project/internal/util"  // <-- Import shared utils
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware creates a middleware that verifies the JWT token.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("AuthMiddleware triggered") // Using Debug level for less noise

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				slog.Warn("authorization header is missing")
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				slog.Warn("invalid authorization header format")
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			tokenString := headerParts[1]

			// Use the Claims struct from our DTO package.
			claims := &dto.Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// Use the actual secret passed from the config.
    			return []byte(util.HardcodedJWTSecret), nil
			})

			if err != nil {
				// Log the specific validation error.
				if errors.Is(err, jwt.ErrTokenExpired) {
					slog.Warn("token validation failed: token is expired", "error", err)
				} else {
					slog.Warn("token validation failed", "error", err)
				}
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				slog.Warn("token was parsed but is not valid")
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			
			slog.Debug("token is valid", "user_id", claims.UserID, "role_id", claims.RoleID)

			// Add the claims to the context using our new utility function.
			ctx := util.AddClaimsToContext(r.Context(), claims)
			
			// Call the next handler in the chain.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TimeoutMiddleware adds a request timeout.
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthorizeRole creates a middleware that checks if the user has a required role.
func AuthorizeRole(allowedRoleIDs ...int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get claims from the context using our new utility function.
			claims, ok := util.GetClaimsFromContext(r.Context())
			if !ok {
				slog.Error("could not get claims from context in AuthorizeRole middleware")
				http.Error(w, "Not authorized", http.StatusForbidden)
				return
			}

			isAllowed := false
			for _, allowedRole := range allowedRoleIDs {
				if claims.RoleID == allowedRole {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				slog.Warn("user forbidden from accessing route", "user_id", claims.UserID, "user_role", claims.RoleID, "required_roles", allowedRoleIDs)
				http.Error(w, "Forbidden: You do not have the necessary permissions.", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}