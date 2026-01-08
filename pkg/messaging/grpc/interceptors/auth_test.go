package interceptors

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestJWTAuthenticator_Authenticate(t *testing.T) {
	auth := NewJWTAuthenticator("test-secret-key", "test-issuer")

	// Generate a valid token
	token, err := auth.GenerateToken("user123", "user@example.com", []string{"admin"}, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Test valid token
	ctx := context.Background()
	authCtx, err := auth.Authenticate(ctx, token)
	if err != nil {
		t.Fatalf("authentication failed: %v", err)
	}

	// Verify user info was added to context
	userID, ok := GetUserID(authCtx)
	if !ok || userID != "user123" {
		t.Errorf("expected user_id=user123, got %s", userID)
	}

	email, ok := GetUserEmail(authCtx)
	if !ok || email != "user@example.com" {
		t.Errorf("expected email=user@example.com, got %s", email)
	}

	roles, ok := GetUserRoles(authCtx)
	if !ok || len(roles) != 1 || roles[0] != "admin" {
		t.Errorf("expected roles=[admin], got %v", roles)
	}
}

func TestJWTAuthenticator_InvalidToken(t *testing.T) {
	auth := NewJWTAuthenticator("test-secret-key", "test-issuer")

	ctx := context.Background()
	_, err := auth.Authenticate(ctx, "invalid-token")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestJWTAuthenticator_ExpiredToken(t *testing.T) {
	auth := NewJWTAuthenticator("test-secret-key", "test-issuer")

	// Generate an expired token
	token, err := auth.GenerateToken("user123", "user@example.com", nil, -time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	ctx := context.Background()
	_, err = auth.Authenticate(ctx, token)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestAPIKeyAuthenticator_Authenticate(t *testing.T) {
	keys := map[string]APIKeyInfo{
		"test-api-key-123": {
			UserID: "user456",
			Email:  "api@example.com",
			Roles:  []string{"service"},
		},
	}

	auth := NewAPIKeyAuthenticator(keys)

	// Test valid API key
	ctx := context.Background()
	authCtx, err := auth.Authenticate(ctx, "test-api-key-123")
	if err != nil {
		t.Fatalf("authentication failed: %v", err)
	}

	// Verify user info
	userID, ok := GetUserID(authCtx)
	if !ok || userID != "user456" {
		t.Errorf("expected user_id=user456, got %s", userID)
	}

	email, ok := GetUserEmail(authCtx)
	if !ok || email != "api@example.com" {
		t.Errorf("expected email=api@example.com, got %s", email)
	}
}

func TestAPIKeyAuthenticator_InvalidKey(t *testing.T) {
	auth := NewAPIKeyAuthenticator(map[string]APIKeyInfo{})

	ctx := context.Background()
	_, err := auth.Authenticate(ctx, "invalid-key")
	if err == nil {
		t.Error("expected error for invalid API key")
	}
}

func TestExtractToken_Bearer(t *testing.T) {
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("authorization", "Bearer test-token"),
	)

	token, err := extractToken(ctx)
	if err != nil {
		t.Fatalf("failed to extract token: %v", err)
	}

	if token != "test-token" {
		t.Errorf("expected token=test-token, got %s", token)
	}
}

func TestExtractToken_APIKey(t *testing.T) {
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("x-api-key", "my-api-key"),
	)

	token, err := extractToken(ctx)
	if err != nil {
		t.Fatalf("failed to extract token: %v", err)
	}

	if token != "my-api-key" {
		t.Errorf("expected token=my-api-key, got %s", token)
	}
}

func TestExtractToken_Missing(t *testing.T) {
	ctx := context.Background()

	_, err := extractToken(ctx)
	if err == nil {
		t.Error("expected error for missing token")
	}
}

func TestShouldSkip(t *testing.T) {
	skipMethods := []string{"/health", "/metrics"}

	if !shouldSkip("/health", skipMethods) {
		t.Error("expected /health to be skipped")
	}

	if shouldSkip("/api/users", skipMethods) {
		t.Error("expected /api/users to not be skipped")
	}
}
