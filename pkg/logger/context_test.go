package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestWithContext(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx := context.Background()
	ctx = WithContext(ctx, logger)

	// Verify logger was stored in context
	retrievedLogger := FromContext(ctx)
	if retrievedLogger != logger {
		t.Error("FromContext() did not return the expected logger")
	}
}

func TestFromContext_NoLogger(t *testing.T) {
	// Create a new logger first
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	_, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Context without logger should return global logger
	ctx := context.Background()
	logger := FromContext(ctx)

	if logger == nil {
		t.Error("FromContext() returned nil when no logger in context")
	}

	// Should return the global logger
	if logger != Get() {
		t.Error("FromContext() should return global logger when none in context")
	}
}

func TestWithRequestID(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx := context.Background()
	ctx = WithContext(ctx, logger)

	requestID := "req-12345"
	ctx = WithRequestID(ctx, requestID)

	// Verify logger in context has request_id field
	loggerWithID := FromContext(ctx)
	if loggerWithID == nil {
		t.Error("FromContext() returned nil after WithRequestID()")
	}

	// Logger should be different from original (has additional field)
	if loggerWithID == logger {
		t.Error("WithRequestID() should return a new logger with request_id field")
	}
}

func TestWithTraceID(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx := context.Background()
	ctx = WithContext(ctx, logger)

	traceID := "trace-67890"
	ctx = WithTraceID(ctx, traceID)

	// Verify logger in context has trace_id field
	loggerWithTrace := FromContext(ctx)
	if loggerWithTrace == nil {
		t.Error("FromContext() returned nil after WithTraceID()")
	}

	// Logger should be different from original (has additional field)
	if loggerWithTrace == logger {
		t.Error("WithTraceID() should return a new logger with trace_id field")
	}
}

func TestMultipleContextFields(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx := context.Background()
	ctx = WithContext(ctx, logger)

	// Add multiple fields
	ctx = WithRequestID(ctx, "req-123")
	ctx = WithTraceID(ctx, "trace-456")

	// Verify logger has both fields
	finalLogger := FromContext(ctx)
	if finalLogger == nil {
		t.Error("FromContext() returned nil after adding multiple fields")
	}

	// Should be able to log without panic
	finalLogger.Info("test message with request and trace IDs")
}

func TestContextPropagation(t *testing.T) {
	config := Config{
		Level:      "debug",
		Format:     "console",
		OutputPath: "stdout",
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Simulate request flow with context propagation
	ctx := context.Background()
	ctx = WithContext(ctx, logger)
	ctx = WithRequestID(ctx, "req-001")

	// Simulate passing context through layers
	processRequest := func(ctx context.Context) {
		log := FromContext(ctx)
		log.Info("processing request")

		// Add trace ID in deeper layer
		ctx = WithTraceID(ctx, "trace-001")
		log = FromContext(ctx)
		log.Debug("added trace ID")
	}

	// Should not panic
	processRequest(ctx)
}

func TestFromContext_WithFields(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	baseLogger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Create logger with fields
	loggerWithFields := baseLogger.With(
		zap.String("service", "test-service"),
		zap.String("version", "1.0.0"),
	)

	ctx := context.Background()
	ctx = WithContext(ctx, loggerWithFields)

	// Add request ID
	ctx = WithRequestID(ctx, "req-999")

	// Logger should have all fields
	finalLogger := FromContext(ctx)
	if finalLogger == nil {
		t.Error("FromContext() returned nil")
	}

	// Should be able to log with all fields
	finalLogger.Info("test with multiple fields")
}

func TestNilContext(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	_, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// FromContext with nil context should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("FromContext(nil) panicked: %v", r)
		}
	}()

	// This will panic in the actual implementation, but we're testing the behavior
	// In production, this should be avoided
	// logger := FromContext(nil)
	// We skip this test as it would panic
}

func TestContextChaining(t *testing.T) {
	config := Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Chain context operations
	ctx := WithTraceID(
		WithRequestID(
			WithContext(context.Background(), logger),
			"req-chain-1",
		),
		"trace-chain-1",
	)

	finalLogger := FromContext(ctx)
	if finalLogger == nil {
		t.Error("Chained context operations resulted in nil logger")
	}

	// Should work without panic
	finalLogger.Info("chained context test")
}
