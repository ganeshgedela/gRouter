package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const loggerKey contextKey = "logger"

// WithContext adds a logger to the context
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves a logger from the context
// If no logger is found, returns the global logger
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger
	}
	return Get()
}

// WithRequestID adds a request ID to the logger in context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	logger := FromContext(ctx).With(zap.String("request_id", requestID))
	return WithContext(ctx, logger)
}

// WithTraceID adds a trace ID to the logger in context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	logger := FromContext(ctx).With(zap.String("trace_id", traceID))
	return WithContext(ctx, logger)
}
