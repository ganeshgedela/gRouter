package nats

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddleware(t *testing.T) {
	core, obs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	mw := LoggingMiddleware(logger)
	handler := mw(func(ctx context.Context, subject string, env *MessageEnvelope) error {
		return nil
	})

	env := &MessageEnvelope{
		ID:   "test-id",
		Type: "test-type",
	}

	err := handler(context.Background(), "test.subject", env)
	assert.NoError(t, err)

	assert.Equal(t, 1, obs.Len())
	assert.Equal(t, "Message processed successfully", obs.All()[0].Message)
}

func TestMetricsMiddleware(t *testing.T) {
	mw := MetricsMiddleware()
	handler := mw(func(ctx context.Context, subject string, env *MessageEnvelope) error {
		return nil
	})

	env := &MessageEnvelope{
		ID:   "test-id",
		Type: "test-type",
	}

	// Reset metrics if possible or just check increment
	before := testutil.ToFloat64(subscribeCounter.WithLabelValues("test.subject", "test-type", "success"))

	err := handler(context.Background(), "test.subject", env)
	assert.NoError(t, err)

	after := testutil.ToFloat64(subscribeCounter.WithLabelValues("test.subject", "test-type", "success"))
	assert.Equal(t, before+1, after)
}

func TestTracingMiddleware(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exporter)))
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("test")

	mw := TracingMiddleware(tracer)
	handler := mw(func(ctx context.Context, subject string, env *MessageEnvelope) error {
		return nil
	})

	env := &MessageEnvelope{
		ID:       "test-id",
		Type:     "test-type",
		Metadata: make(map[string]string),
	}

	err := handler(context.Background(), "test.subject", env)
	assert.NoError(t, err)

	spans := exporter.GetSpans()
	assert.Len(t, spans, 1)
	assert.Equal(t, "messaging.receive test.subject", spans[0].Name)
}

func TestPublisherTracingMiddleware(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exporter)))
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("test")

	mw := PublisherTracingMiddleware(tracer)
	publishFunc := mw(func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
		return nil
	})

	err := publishFunc(context.Background(), "test.subject", "test-type", nil, nil)
	assert.NoError(t, err)

	spans := exporter.GetSpans()
	assert.Len(t, spans, 1)
	assert.Equal(t, "messaging.send test.subject", spans[0].Name)
}
