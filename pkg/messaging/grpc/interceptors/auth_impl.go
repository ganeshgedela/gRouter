package interceptors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthenticator implements JWT-based authentication
type JWTAuthenticator struct {
	secretKey []byte
	issuer    string
}

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// NewJWTAuthenticator creates a new JWT authenticator
func NewJWTAuthenticator(secretKey string, issuer string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// Authenticate validates a JWT token and adds user info to context
func (a *JWTAuthenticator) Authenticate(ctx context.Context, tokenString string) (context.Context, error) {
	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secretKey, nil
	})

	if err != nil {
		return ctx, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return ctx, fmt.Errorf("invalid token claims")
	}

	// Verify issuer if configured
	if a.issuer != "" && claims.Issuer != a.issuer {
		return ctx, fmt.Errorf("invalid issuer")
	}

	// Verify expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return ctx, fmt.Errorf("token expired")
	}

	// Add user info to context
	ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
	if claims.Email != "" {
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
	}
	if len(claims.Roles) > 0 {
		ctx = context.WithValue(ctx, UserRolesKey, claims.Roles)
	}

	return ctx, nil
}

// GenerateToken generates a JWT token (useful for testing)
func (a *JWTAuthenticator) GenerateToken(userID, email string, roles []string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    a.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secretKey)
}

// APIKeyAuthenticator implements API key-based authentication
type APIKeyAuthenticator struct {
	keys map[string]APIKeyInfo
}

// APIKeyInfo holds information about an API key
type APIKeyInfo struct {
	UserID string
	Email  string
	Roles  []string
}

// NewAPIKeyAuthenticator creates a new API key authenticator
func NewAPIKeyAuthenticator(keys map[string]APIKeyInfo) *APIKeyAuthenticator {
	return &APIKeyAuthenticator{
		keys: keys,
	}
}

// Authenticate validates an API key and adds user info to context
func (a *APIKeyAuthenticator) Authenticate(ctx context.Context, apiKey string) (context.Context, error) {
	// Trim whitespace
	apiKey = strings.TrimSpace(apiKey)

	// Look up the key
	keyInfo, exists := a.keys[apiKey]
	if !exists {
		return ctx, fmt.Errorf("invalid API key")
	}

	// Add user info to context
	ctx = context.WithValue(ctx, UserIDKey, keyInfo.UserID)
	if keyInfo.Email != "" {
		ctx = context.WithValue(ctx, UserEmailKey, keyInfo.Email)
	}
	if len(keyInfo.Roles) > 0 {
		ctx = context.WithValue(ctx, UserRolesKey, keyInfo.Roles)
	}

	return ctx, nil
}

// AddKey adds a new API key
func (a *APIKeyAuthenticator) AddKey(apiKey string, info APIKeyInfo) {
	a.keys[apiKey] = info
}

// RemoveKey removes an API key
func (a *APIKeyAuthenticator) RemoveKey(apiKey string) {
	delete(a.keys, apiKey)
}
