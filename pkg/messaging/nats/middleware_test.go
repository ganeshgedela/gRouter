package nats

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddleware(t *testing.T) {
	core, obs := observer.New(zap.DebugLevel)
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
	before := testutil.ToFloat64(subscribeCounter.WithLabelValues("test.subject", "test-type", "success", "unknown"))

	err := handler(context.Background(), "test.subject", env)
	assert.NoError(t, err)

	after := testutil.ToFloat64(subscribeCounter.WithLabelValues("test.subject", "test-type", "success", "unknown"))
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
	assert.Equal(t, "test.subject process", spans[0].Name)
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
	assert.Equal(t, "test.subject send", spans[0].Name)
}

func TestTracingHelpers(t *testing.T) {
	// Test context propagation carrier
	md := make(map[string]string)
	carrier := metadataCarrier(md)

	// Set/Get
	carrier.Set("key1", "value1")
	val := carrier.Get("key1")
	assert.Equal(t, "value1", val)

	val = carrier.Get("nonexistent")
	assert.Equal(t, "", val)

	// Keys
	keys := carrier.Keys()
	assert.Contains(t, keys, "key1")
}

func TestRecoveryMiddleware(t *testing.T) {
	core, obs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)

	mw := RecoveryMiddleware(logger)
	handler := mw(func(ctx context.Context, subject string, env *MessageEnvelope) error {
		panic("simulated panic")
	})

	err := handler(context.Background(), "test.subject", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic recovered: simulated panic")

	assert.Equal(t, "Panic recovered in subscriber", obs.All()[0].Message)
}

func TestCircuitBreakerMiddleware(t *testing.T) {
	cfg := CircuitBreakerConfig{
		Enabled:       true,
		MaxRequests:   1,
		Interval:      time.Second * 10,
		Timeout:       time.Second * 10,
		TripThreshold: 2,
	}
	logger := zap.NewNop()
	cb := NewCircuitBreaker(cfg, logger)
	mw := CircuitBreakerMiddleware(cb, logger)

	failHandler := func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
		return fmt.Errorf("simulated failure")
	}

	successHandler := func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
		return nil
	}

	// 1. Success should pass
	wrappedSuccess := mw(successHandler)
	err := wrappedSuccess(context.Background(), "test.cb", "type", nil, nil)
	assert.NoError(t, err)

	// 2. Failures to trip
	wrappedFail := mw(failHandler)
	// Failure 1
	err = wrappedFail(context.Background(), "test.cb", "type", nil, nil)
	assert.Error(t, err)
	assert.Equal(t, "simulated failure", err.Error())

	// Failure 2 (Should trip because threshold is 2, but gobreaker trips AFTER the threshold is reached or exceeded?
	// Our logic: counts.ConsecutiveFailures >= threshold (2).
	// After 1 failure, Consecutive = 1.
	// After 2 failures, Consecutive = 2. 2 >= 2 -> Trip.
	err = wrappedFail(context.Background(), "test.cb", "type", nil, nil)
	assert.Error(t, err)

	// 3. Next call should be rejected by CB
	err = wrappedSuccess(context.Background(), "test.cb", "type", nil, nil)
	assert.Error(t, err)
	assert.Equal(t, "service temporarily unavailable", err.Error())
}

func TestRetryMiddleware(t *testing.T) {
	cfg := RetryConfig{
		Enabled:         true,
		MaxAttempts:     3,
		InitialInterval: time.Millisecond * 10,
		Multiplier:      2.0,
		MaxInterval:     time.Millisecond * 100,
	}
	mw := RetryMiddleware(cfg)

	// 1. Success immediately
	calls := 0
	successHandler := func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
		calls++
		return nil
	}
	err := mw(successHandler)(context.Background(), "test.retry", "type", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, calls)

	// 2. Success after retry (Fail once, then succeed)
	calls = 0
	flakeyHandler := func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
		calls++
		if calls < 2 {
			return fmt.Errorf("transient error")
		}
		return nil
	}
	err = mw(flakeyHandler)(context.Background(), "test.retry", "type", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, calls)

	// 3. Fail after max retries (Fail 3 times)
	calls = 0
	failHandler := func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
		calls++
		return fmt.Errorf("persistent error")
	}
	err = mw(failHandler)(context.Background(), "test.retry", "type", nil, nil)
	assert.Error(t, err)
	assert.Equal(t, "persistent error", err.Error())
	assert.Equal(t, 3, calls)
}

func TestTimeoutMiddleware(t *testing.T) {
	cfg := TimeoutConfig{
		Enabled: true,
		Default: time.Millisecond * 50,
	}
	mw := TimeoutMiddleware(cfg)

	// 1. Verify deadline set
	checkDeadline := func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
		_, ok := ctx.Deadline()
		assert.True(t, ok, "deadline should be set")
		return nil
	}
	err := mw(checkDeadline)(context.Background(), "test.timeout", "type", nil, nil)
	assert.NoError(t, err)

	// 2. Verify timeout expiration
	// Wrapper usually doesn't return error itself unless next function respects context.
	// So we simulate a next function that waits on context or timer.
	sleepHandler := func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
		select {
		case <-time.After(time.Millisecond * 100):
			return fmt.Errorf("finished normally, should have timed out")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	err = mw(sleepHandler)(context.Background(), "test.timeout", "type", nil, nil)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestRateLimitMiddleware(t *testing.T) {
	cfg := RateLimitConfig{
		Enabled:           true,
		RequestsPerSecond: 1, // 1 request / sec
		Burst:             1,
	}
	limiter := NewRateLimiter(cfg)
	mw := RateLimitMiddleware(limiter)
	noop := func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
		return nil
	}

	// 1. First call passes (within burst)
	err := mw(noop)(context.Background(), "test.ratelimit", "type", nil, nil)
	assert.NoError(t, err)

	// 2. Second call immediately should block ~1s.
	// We use a short timeout to enforce error.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	err = mw(noop)(ctx, "test.ratelimit", "type", nil, nil)
	assert.Error(t, err)
	// rate.Limiter Wait returns a specific error if it knows it will exceed deadline
	assert.Contains(t, err.Error(), "exceed context deadline")
}

type mockValidator struct {
	err error
}

func (m *mockValidator) Validate(msgType string, data []byte) error {
	return m.err
}

func TestValidatorMiddleware(t *testing.T) {
	p := &NATSPublisher{}
	mw := ValidatorMiddleware(p)
	noop := func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
		return nil
	}

	// 1. No validator set -> should pass
	err := mw(noop)(context.Background(), "test.val", "type", nil, nil)
	assert.NoError(t, err)

	// 2. Validator set -> Pass
	p.SetValidator(&mockValidator{})
	err = mw(noop)(context.Background(), "test.val", "type", nil, nil)
	assert.NoError(t, err)

	// 3. Validator set -> Fail
	p.SetValidator(&mockValidator{err: fmt.Errorf("bad data")})
	err = mw(noop)(context.Background(), "test.val", "type", nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// 4. Test types
	// 4a. []byte
	err = mw(noop)(context.Background(), "test.val", "type", []byte("data"), nil)
	assert.Error(t, err)

	// 4b. interface{} (map)
	p.SetValidator(&mockValidator{})
	err = mw(noop)(context.Background(), "test.val", "type", map[string]string{"foo": "bar"}, nil)
	assert.NoError(t, err)
}
