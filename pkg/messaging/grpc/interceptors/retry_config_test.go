package interceptors

import (
	"testing"
)

// TestDefaultRetryConfig tests the default retry configuration.
func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts 3, got %d", config.MaxAttempts)
	}

	if len(config.RetryableCodes) == 0 {
		t.Error("expected retryable codes to be configured")
	}
}
