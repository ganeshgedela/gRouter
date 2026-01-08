package interceptors

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Authenticator defines the interface for authentication
type Authenticator interface {
	Authenticate(ctx context.Context, token string) (context.Context, error)
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Enabled       bool
	Authenticator Authenticator
	SkipMethods   []string // Methods to skip authentication
}

// contextKey is a private type for context keys
type contextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey contextKey = "user_id"
	// UserEmailKey is the context key for user email
	UserEmailKey contextKey = "user_email"
	// UserRolesKey is the context key for user roles
	UserRolesKey contextKey = "user_roles"
)

// AuthUnaryInterceptor returns a unary server interceptor for authentication
func AuthUnaryInterceptor(config AuthConfig) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !config.Enabled || config.Authenticator == nil {
			return handler(ctx, req)
		}

		// Check if method should skip auth
		if shouldSkip(info.FullMethod, config.SkipMethods) {
			return handler(ctx, req)
		}

		// Extract token from metadata
		token, err := extractToken(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
		}

		// Authenticate
		authCtx, err := config.Authenticator.Authenticate(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return handler(authCtx, req)
	}
}

// AuthStreamInterceptor returns a stream server interceptor for authentication
func AuthStreamInterceptor(config AuthConfig) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if !config.Enabled || config.Authenticator == nil {
			return handler(srv, ss)
		}

		// Check if method should skip auth
		if shouldSkip(info.FullMethod, config.SkipMethods) {
			return handler(srv, ss)
		}

		// Extract token from metadata
		ctx := ss.Context()
		token, err := extractToken(ctx)
		if err != nil {
			return status.Error(codes.Unauthenticated, "missing or invalid token")
		}

		// Authenticate
		authCtx, err := config.Authenticator.Authenticate(ctx, token)
		if err != nil {
			return status.Error(codes.Unauthenticated, err.Error())
		}

		// Wrap the stream with authenticated context
		wrappedStream := &authenticatedStream{
			ServerStream: ss,
			ctx:          authCtx,
		}

		return handler(srv, wrappedStream)
	}
}

// authenticatedStream wraps a grpc.ServerStream with an authenticated context
type authenticatedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *authenticatedStream) Context() context.Context {
	return s.ctx
}

// extractToken extracts the authentication token from gRPC metadata
func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing metadata")
	}

	// Try "authorization" header first (common for JWT)
	authHeaders := md.Get("authorization")
	if len(authHeaders) > 0 {
		// Expected format: "Bearer <token>"
		parts := strings.SplitN(authHeaders[0], " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1], nil
		}
		// Also accept bare token
		return authHeaders[0], nil
	}

	// Try "x-api-key" header (common for API keys)
	apiKeyHeaders := md.Get("x-api-key")
	if len(apiKeyHeaders) > 0 {
		return apiKeyHeaders[0], nil
	}

	return "", fmt.Errorf("no authentication token found")
}

// shouldSkip checks if a method should skip authentication
func shouldSkip(method string, skipMethods []string) bool {
	for _, skip := range skipMethods {
		if method == skip {
			return true
		}
	}
	return false
}

// GetUserID retrieves the user ID from context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetUserEmail retrieves the user email from context
func GetUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

// GetUserRoles retrieves the user roles from context
func GetUserRoles(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(UserRolesKey).([]string)
	return roles, ok
}
