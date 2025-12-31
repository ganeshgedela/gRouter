package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware that verifies OIDC ID tokens
func AuthMiddleware(cfg AuthConfig) gin.HandlerFunc {
	// If disabled, just pass through
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	provider, err := oidc.NewProvider(context.Background(), cfg.Issuer)
	if err != nil {
		// If provider initialization fails, we panic because auth is critical but misconfigured
		// In production, might want to retry or error out gracefully at startup.
		// For middleware factory, we usually return error or panic.
		// Since gin.HandlerFunc signature doesn't allow error return, we'll log and panic.
		panic(fmt.Sprintf("failed to init OIDC provider: %v", err))
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.Audience,
	})

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		// Verify token
		idToken, err := verifier.Verify(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			return
		}

		// Store claims/token in context
		c.Set("token", idToken)

		// Extract claims if needed
		var claims struct {
			Email    string `json:"email"`
			Verified bool   `json:"email_verified"`
			Sub      string `json:"sub"`
		}
		if err := idToken.Claims(&claims); err == nil {
			c.Set("user_email", claims.Email)
			c.Set("user_id", claims.Sub)
		}

		c.Next()
	}
}
