// File: internal/util/context.go
package util

import (
    "context"
    "crm-project/internal/dto"
)

// Define a new type for our context key to prevent collisions.
type contextKey string
const userContextKey = contextKey("user")

// GetClaimsFromContext now lives in the neutral 'util' package.
func GetClaimsFromContext(ctx context.Context) (*dto.Claims, bool) {
    claims, ok := ctx.Value(userContextKey).(*dto.Claims)
    return claims, ok
}

// We also need a way to add the claims to the context.
func AddClaimsToContext(ctx context.Context, claims *dto.Claims) context.Context {
    return context.WithValue(ctx, userContextKey, claims)
}


const HardcodedJWTSecret = "my_super_secret_debug_key_12345"
