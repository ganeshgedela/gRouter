package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"grouter/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestDefaultSecurityHeaders(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	cfg := config.WebConfig{
		Port:    18091,
		Mode:    "test",
		Logging: config.LoggingConfig{Enabled: false},
		// Security config disabled
		Security: config.SecurityConfig{Enabled: false},
	}

	srv, err := NewServer("test-headers", "Test", cfg, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Add a test route
	srv.Engine().GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	srv.Engine().ServeHTTP(w, req)

	// Check default security headers are present
	headers := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
		"X-Server-Version":       "1.0.0",
	}

	for header, expected := range headers {
		actual := w.Header().Get(header)
		if actual != expected {
			t.Errorf("%s: expected %q, got %q", header, expected, actual)
		}
	}
}

func TestErrorResponseIncludesRequestID(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	gin.SetMode(gin.TestMode)

	cfg := config.WebConfig{
		Port:    18092,
		Mode:    "test",
		Logging: config.LoggingConfig{Enabled: false},
	}

	srv, _ := NewServer("test-error-id", "Test", cfg, logger)

	// Route that returns error
	srv.Engine().GET("/error", func(c *gin.Context) {
		RespondError(c, ErrBadRequest)
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()

	srv.Engine().ServeHTTP(w, req)

	// Check request ID in response
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("X-Request-ID header not set")
	}

	// Check request ID in error body
	body := w.Body.String()
	if body == "" {
		t.Error("Empty response body")
	}

	// Should contain request_id field
	if !contains(body, "request_id") {
		t.Error("Error response missing request_id field")
	}

	if !contains(body, requestID) {
		t.Error("Error response request_id doesn't match header")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) && hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
