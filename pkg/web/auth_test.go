package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestJWTAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret-key"
	cfg := JWTConfig{
		Enabled:   true,
		Secret:    secret,
		Algorithm: "HS256",
	}

	// Generate test token
	token, err := GenerateJWT("user123", "test@example.com", []string{"admin"}, secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
	}{
		{
			name:         "valid token",
			authHeader:   "Bearer " + token,
			expectedCode: http.StatusOK,
		},
		{
			name:         "missing token",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid format",
			authHeader:   "InvalidFormat",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid token",
			authHeader:   "Bearer invalid.token.here",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(JWTAuthMiddleware(cfg))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestAPIKeyAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	store := NewAPIKeyStore()
	store.Add("valid-key", APIKeyInfo{
		UserID: "user123",
		Email:  "test@example.com",
		Roles:  []string{"user"},
		Active: true,
	})
	store.Add("inactive-key", APIKeyInfo{
		UserID: "user456",
		Email:  "inactive@example.com",
		Active: false,
	})

	tests := []struct {
		name         string
		apiKey       string
		expectedCode int
	}{
		{
			name:         "valid API key",
			apiKey:       "valid-key",
			expectedCode: http.StatusOK,
		},
		{
			name:         "missing API key",
			apiKey:       "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid API key",
			apiKey:       "invalid-key",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inactive API key",
			apiKey:       "inactive-key",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(APIKeyAuthMiddleware(store, "X-API-Key"))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestRBACMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		userRoles     []string
		requiredRoles []string
		expectedCode  int
	}{
		{
			name:          "user has required role",
			userRoles:     []string{"admin", "user"},
			requiredRoles: []string{"admin"},
			expectedCode:  http.StatusOK,
		},
		{
			name:          "user missing required role",
			userRoles:     []string{"user"},
			requiredRoles: []string{"admin"},
			expectedCode:  http.StatusForbidden,
		},
		{
			name:          "user has one of multiple required roles",
			userRoles:     []string{"editor"},
			requiredRoles: []string{"admin", "editor"},
			expectedCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("user_roles", tt.userRoles)
				c.Next()
			})
			router.Use(RBACMiddleware(nil, tt.requiredRoles...))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}
