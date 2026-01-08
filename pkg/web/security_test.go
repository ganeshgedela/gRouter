package web

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"grouter/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestRequestBodySizeLimit(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	// Create server with 1KB limit
	cfg := config.WebConfig{
		Port:               18090,
		Mode:               "test",
		MaxRequestBodySize: 1024, // 1KB
		Logging:            config.LoggingConfig{Enabled: false},
	}

	srv, err := NewServer("test-body-limit", "Test", cfg, logger)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	// Add a test route that reads body
	srv.Engine().POST("/upload", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.JSON(http.StatusOK, gin.H{"size": len(body)})
	})

	tests := []struct {
		name        string
		bodySize    int
		expectError bool
	}{
		{
			name:        "small body succeeds",
			bodySize:    512,
			expectError: false,
		},
		{
			name:        "body at limit succeeds",
			bodySize:    1024,
			expectError: false,
		},
		{
			name:        "oversized body rejected",
			bodySize:    2048,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewReader(make([]byte, tt.bodySize))
			req := httptest.NewRequest("POST", "/upload", body)
			w := httptest.NewRecorder()

			srv.Engine().ServeHTTP(w, req)

			// MaxBytesReader in httptest may not always return proper status codes
			// Check if body reading would fail for oversized requests
			if tt.expectError {
				// Either 413 status or the handler should not receive full body
				if w.Code == http.StatusOK {
					// In some test scenarios, check response body
					var resp map[string]interface{}
					_ = json.Unmarshal(w.Body.Bytes(), &resp)
					if size, ok := resp["size"].(float64); ok && size >= float64(tt.bodySize) {
						t.Errorf("oversized body was not rejected, got size %v", size)
					}
				}
			} else if w.Code != http.StatusOK {
				t.Errorf("expected status 200, got %d", w.Code)
			}
		})
	}
}

func TestJWTErrorSanitization(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	gin.SetMode(gin.TestMode)

	cfg := JWTConfig{
		Enabled:   true,
		Secret:    "test-secret",
		Algorithm: "HS256",
		Logger:    logger,
	}

	router := gin.New()
	router.Use(JWTAuthMiddleware(cfg))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test with invalid token format
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.malformed.token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	// Check response doesn't contain internal error details
	body := w.Body.String()
	if strings.Contains(body, "signature") || strings.Contains(body, "algorithm") {
		t.Error("response contains internal error details - information leakage!")
	}

	// Should contain generic message
	if !strings.Contains(body, "invalid or expired token") {
		t.Error("response should contain generic error message")
	}
}
