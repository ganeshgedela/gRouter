package web

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// APIKeyConfig holds API key authentication configuration
type APIKeyConfig struct {
	Enabled bool
	Header  string // Header name to check (default: X-API-Key)
	Keys    map[string]APIKeyInfo
}

// APIKeyInfo contains information about an API key
type APIKeyInfo struct {
	UserID string
	Email  string
	Roles  []string
	Active bool
}

// APIKeyStore manages API keys (thread-safe)
type APIKeyStore struct {
	mu   sync.RWMutex
	keys map[string]APIKeyInfo
}

// NewAPIKeyStore creates a new API key store
func NewAPIKeyStore() *APIKeyStore {
	return &APIKeyStore{
		keys: make(map[string]APIKeyInfo),
	}
}

// Add adds an API key
func (s *APIKeyStore) Add(key string, info APIKeyInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keys[key] = info
}

// Get retrieves API key info
func (s *APIKeyStore) Get(key string) (APIKeyInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, exists := s.keys[key]
	return info, exists
}

// Remove removes an API key
func (s *APIKeyStore) Remove(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.keys, key)
}

// APIKeyAuthMiddleware validates API keys
func APIKeyAuthMiddleware(store *APIKeyStore, headerName string) gin.HandlerFunc {
	if headerName == "" {
		headerName = "X-API-Key"
	}

	return func(c *gin.Context) {
		apiKey := c.GetHeader(headerName)
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "MISSING_API_KEY",
				Message: "API key is required",
				Code:    "UNAUTHORIZED",
			})
			return
		}

		// Validate API key
		keyInfo, exists := store.Get(apiKey)
		if !exists || !keyInfo.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "INVALID_API_KEY",
				Message: "invalid or inactive API key",
				Code:    "UNAUTHORIZED",
			})
			return
		}

		// Store user info in context
		c.Set("user_id", keyInfo.UserID)
		c.Set("user_email", keyInfo.Email)
		c.Set("user_roles", keyInfo.Roles)
		c.Set("auth_method", "apikey")

		c.Next()
	}
}
