# NATS Middleware Development Guide

## Overview

Middleware in the NATS package provides a powerful way to add cross-cutting concerns to your message publishing and subscribing operations. This guide covers the available middleware, how to use them, and how to create custom middleware.

## Available Middleware

### 1. Recovery Middleware ⛑️

**Purpose**: Catches panics in handlers or publishing operations and converts them to errors.

**Usage**:
```go
publisher.Use(PublisherRecoveryMiddleware(logger))
subscriber.Use(RecoveryMiddleware(logger))
```

**Behavior**:
- Catches any panic in the handler chain
- Logs the panic with stack trace
- Returns an error instead of crashing

**Position**: Should be **outermost** (first registered) to catch all panics.

---

### 2. Metrics Middleware 📊

**Purpose**: Emits Prometheus metrics for observability.

**Metrics Exported**:
- `nats_messages_total{service, subject, status}` - Total message count
- `nats_message_duration_seconds{service, subject}` - Message processing duration
- `nats_message_errors_total{service, subject}` - Error count

**Usage**:
```go
publisher.Use(PublisherMetricsMiddleware())
subscriber.Use(MetricsMiddleware())
```

**Configuration**:
```go
cfg.Middleware.Metrics.Enabled = true
cfg.Middleware.Metrics.Path = "/metrics"  // For HTTP metrics endpoint
```

---

### 3. Logging Middleware 📝

**Purpose**: Structured logging of all message operations using Zap.

**Log Fields**:
- `subject` - NATS subject
- `message_id` - Envelope ID
- `message_type` - Message type
- `duration` - Processing time
- `error` - Error if any

**Usage**:
```go
publisher.Use(PublisherLoggingMiddleware(logger))
subscriber.Use(LoggingMiddleware(logger))
```

**Configuration**:
```go
cfg.Middleware.Logging.Enabled = true
```

---

### 4. Tracing Middleware 🔍

**Purpose**: OpenTelemetry distributed tracing integration.

**Features**:
- Creates spans for each message operation
- Propagates trace context via metadata
- Links request-reply operations

**Usage**:
```go
tracer := otel.Tracer("nats")
publisher.Use(PublisherTracingMiddleware(tracer))
subscriber.Use(TracingMiddleware(tracer))
```

**Configuration**:
```go
cfg.Middleware.Tracing.Enabled = true
```

---

### 5. Timeout Middleware ⏱️

**Purpose**: Enforces operation timeouts.

**Usage**:
```go
publisher.Use(TimeoutMiddleware(cfg.Middleware.Timeout))
```

**Configuration**:
```go
cfg.Middleware.Timeout.Enabled = true
cfg.Middleware.Timeout.Default = 5 * time.Second
```

**Behavior**:
- Creates a context with timeout if none exists
- Cancels operation if timeout is exceeded
- Returns `context.DeadlineExceeded` error

---

### 6. Rate Limit Middleware 🚦

**Purpose**: Token bucket rate limiting for publishers.

**Usage**:
```go
limiter := NewRateLimiter(cfg.Middleware.RateLimit)
publisher.Use(RateLimitMiddleware(limiter))
```

**Configuration**:
```go
cfg.Middleware.RateLimit.Enabled = true
cfg.Middleware.RateLimit.RequestsPerSecond = 100
cfg.Middleware.RateLimit.Burst = 10
```

**Behavior**:
- Allows `RequestsPerSecond` sustained rate
- Allows bursts up to `Burst` size
- Blocks until token is available

---

### 7. Validator Middleware ✅

**Purpose**: Schema validation of message payloads.

**Usage**:
```go
// Register schemas
validator := NewMapValidator()
validator.Register("OrderCreated", func(data []byte) error {
    var order Order
    return json.Unmarshal(data, &order)
})

publisher.SetValidator(validator)
publisher.Use(ValidatorMiddleware(publisher))
```

**Configuration**:
- Validator is always enabled if set
- Thread-safe with `sync.RWMutex`

**Behavior**:
- Validates outgoing messages match registered schema
- Validates incoming messages before handler
- Returns sanitized error on validation failure

---

### 8. Circuit Breaker Middleware ⚡

**Purpose**: Fails fast when error threshold is exceeded.

**Usage**:
```go
cb := NewCircuitBreaker(cfg.Middleware.CircuitBreaker, logger)
publisher.Use(CircuitBreakerMiddleware(cb, logger))
```

**Configuration**:
```go
cfg.Middleware.CircuitBreaker.Enabled = true
cfg.Middleware.CircuitBreaker.TripThreshold = 5      // Failures to trip
cfg.Middleware.CircuitBreaker.Timeout = 60 * time.Second
cfg.Middleware.CircuitBreaker.MaxRequests = 1        // Test requests in half-open
```

**States**:
- **Closed**: Normal operation
- **Open**: Rejecting all requests (after threshold failures)
- **Half-Open**: Testing with limited requests

**Behavior**:
- Counts consecutive failures
- Opens circuit after `TripThreshold` failures
- Returns "service temporarily unavailable" when open
- Logs state transitions

---

### 9. Retry Middleware 🔄

**Purpose**: Exponential backoff retry for transient failures.

**Usage**:
```go
publisher.Use(RetryMiddleware(cfg.Middleware.Retry))
```

**Configuration**:
```go
cfg.Middleware.Retry.Enabled = true
cfg.Middleware.Retry.MaxAttempts = 3
cfg.Middleware.Retry.InitialInterval = 100 * time.Millisecond
cfg.Middleware.Retry.Multiplier = 2.0
cfg.Middleware.Retry.MaxInterval = 5 * time.Second
```

**Backoff Formula**:
```
delay = min(initialInterval * multiplier^(attempt-1), maxInterval)
```

**Behavior**:
- Retries on any error
- Respects context cancellation
- Logs retry attempts

---

## Middleware Ordering

Middleware are executed in **reverse order** of registration (last-in, first-out).

### Recommended Order

```go
// 1. Recovery (catches all panics)
publisher.Use(PublisherRecoveryMiddleware(logger))

// 2. Metrics (measures everything)
publisher.Use(PublisherMetricsMiddleware())

// 3. Logging (logs all operations)
publisher.Use(PublisherLoggingMiddleware(logger))

// 4. Tracing (propagates context)
publisher.Use(PublisherTracingMiddle(tracer))

// 5. Timeout (enforces deadlines)
publisher.Use(TimeoutMiddleware(cfg.Middleware.Timeout))

// 6. Rate Limit (controls traffic)
publisher.Use(RateLimitMiddleware(limiter))

// 7. Validator (checks schema)
publisher.Use(ValidatorMiddleware(publisher))

// 8. Circuit Breaker (fails fast)
publisher.Use(CircuitBreakerMiddleware(cb, logger))

// 9. Retry (retries failures) - Innermost
publisher.Use(RetryMiddleware(cfg.Middleware.Retry))
```

### Execution Flow

```
Request
  ↓
Recovery
  ↓
Metrics (start timer)
  ↓
Logging (log start)
  ↓
Tracing (create span)
  ↓
Timeout (enforce deadline)
  ↓
Rate Limit (wait for token)
  ↓
Validator (check schema)
  ↓
Circuit Breaker (check state)
  ↓
Retry (attempt with backoff)
  ↓
Core Operation (NATS publish/subscribe)
  ↓
Retry (on error)
  ↓
Circuit Breaker (record result)
  ↓
Validator (validate)
  ↓
Rate Limit (return token)
  ↓
Timeout (check deadline)
  ↓
Tracing (end span)
  ↓
Logging (log result)
  ↓
Metrics (record duration)
  ↓
Recovery (catch panic)
  ↓
Response
```

---

## Custom Middleware Development

### Publisher Middleware

```go
type PublisherMiddleware func(PublisherFunc) PublisherFunc

type PublisherFunc func(ctx context.Context, subject, msgType string, 
    data interface{}, opts *PublishOptions) error
```

**Example: Custom Logging Middleware**

```go
func CustomLoggingMiddleware(logger *zap.Logger) PublisherMiddleware {
    return func(next PublisherFunc) PublisherFunc {
        return func(ctx context.Context, subject, msgType string, 
            data interface{}, opts *PublishOptions) error {
            
            start := time.Now()
            
            // Before
            logger.Info("publishing message",
                zap.String("subject", subject),
                zap.String("type", msgType),
            )
            
            // Execute
            err := next(ctx, subject, msgType, data, opts)
            
            // After
            logger.Info("publish complete",
                zap.String("subject", subject),
                zap.Duration("duration", time.Since(start)),
                zap.Error(err),
            )
            
            return err
        }
    }
}

// Use it
publisher.Use(CustomLoggingMiddleware(logger))
```

---

### Subscriber Middleware

```go
type SubscriberMiddleware func(HandlerFunc) HandlerFunc

type HandlerFunc func(ctx context.Context, subject string, 
    env *MessageEnvelope) error
```

**Example: Custom Authorization Middleware**

```go
func AuthorizationMiddleware(requiredRole string) SubscriberMiddleware {
    return func(next HandlerFunc) HandlerFunc {
        return func(ctx context.Context, subject string, env *MessageEnvelope) error {
            // Extract role from metadata
            role, ok := env.Metadata["role"]
            if !ok {
                return fmt.Errorf("missing role in metadata")
            }
            
            // Check authorization
            if role != requiredRole {
                return fmt.Errorf("unauthorized: required role %s", requiredRole)
            }
            
            // Authorized, proceed
            return next(ctx, subject, env)
        }
    }
}

// Use it
subscriber.Use(AuthorizationMiddleware("admin"))
```

---

### Request Middleware

For request-reply patterns:

```go
type RequestMiddleware func(RequestFunc) RequestFunc

type RequestFunc func(ctx context.Context, subject, msgType string, 
    data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error)
```

**Example: Custom Caching Middleware**

```go
func CachingMiddleware(cache *Cache) RequestMiddleware {
    return func(next RequestFunc) RequestFunc {
        return func(ctx context.Context, subject, msgType string, 
            data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error) {
            
            // Generate cache key
            key := fmt.Sprintf("%s:%s:%v", subject, msgType, data)
            
            // Check cache
            if cached, found := cache.Get(key); found {
                return cached.(*MessageEnvelope), nil
            }
            
            // Cache miss, execute request
            response, err := next(ctx, subject, msgType, data, timeout, opts)
            if err != nil {
                return nil, err
            }
            
            // Cache the response
            cache.Set(key, response, 5*time.Minute)
            
            return response, nil
        }
    }
}

// Use it
publisher.UseRequest(CachingMiddleware(cache))
```

---

## Best Practices

### 1. Keep Middleware Focused
Each middleware should do one thing well:
```go
// Good: Single responsibility
func MetricsMiddleware() PublisherMiddleware { ... }
func LoggingMiddleware() PublisherMiddleware { ... }

// Bad: Multiple responsibilities
func MetricsAndLoggingMiddleware() PublisherMiddleware { ... }
```

### 2. Respect Context
Always check context cancellation:
```go
func CustomMiddleware() PublisherMiddleware {
    return func(next PublisherFunc) PublisherFunc {
        return func(ctx context.Context, ...) error {
            // Check context before long operations
            select {
            case <-ctx.Done():
                return ctx.Err()
            default:
            }
            
            return next(ctx, ...)
        }
    }
}
```

### 3. Handle Errors Gracefully
```go
func ResilientMiddleware() PublisherMiddleware {
    return func(next PublisherFunc) PublisherFunc {
        return func(ctx context.Context, ...) error {
            err := next(ctx, ...)
            
            // Log but don't fail on middleware errors
            if internalErr := doSomething(); internalErr != nil {
                logger.Warn("middleware error", zap.Error(internalErr))
            }
            
            // Return original error
            return err
        }
    }
}
```

### 4. Avoid Blocking
```go
// Bad: Blocking on send
func BadMiddleware() PublisherMiddleware {
    return func(next PublisherFunc) PublisherFunc {
        return func(ctx context.Context, ...) error {
            metricsChan <- metric // Blocks if channel is full!
            return next(ctx, ...)
        }
    }
}

// Good: Non-blocking
func GoodMiddleware() PublisherMiddleware {
    return func(next PublisherFunc) PublisherFunc {
        return func(ctx context.Context, ...) error {
            select {
            case metricsChan <- metric:
            default:
                // Drop metric if channel is full
            }
            return next(ctx, ...)
        }
    }
}
```

### 5. Provide Configuration
```go
type CustomMiddlewareConfig struct {
    Enabled   bool
    Threshold int
    Timeout   time.Duration
}

func CustomMiddleware(cfg CustomMiddlewareConfig) PublisherMiddleware {
    if !cfg.Enabled {
        return func(next PublisherFunc) PublisherFunc {
            return next  // No-op if disabled
        }
    }
    
    return func(next PublisherFunc) PublisherFunc {
        return func(ctx context.Context, ...) error {
            // Use cfg.Threshold, cfg.Timeout
            return next(ctx, ...)
        }
    }
}
```

---

## Testing Middleware

### Unit Testing Pattern

```go
func TestCustomMiddleware(t *testing.T) {
    var capturedSubject string
    
    // Mock the next handler
    mockNext := func(ctx context.Context, subject, msgType string, 
        data interface{}, opts *PublishOptions) error {
        capturedSubject = subject
        return nil
    }
    
    // Wrap with middleware
    mw := CustomMiddleware()
    wrapped := mw(mockNext)
    
    // Execute
    err := wrapped(context.Background(), "test.subject", "TestType", nil, nil)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test.subject", capturedSubject)
}
```

---

## Common Patterns

### Conditional Middleware

```go
func ConditionalMiddleware(condition func() bool) PublisherMiddleware {
    return func(next PublisherFunc) PublisherFunc {
        return func(ctx context.Context, ...) error {
            if condition() {
                // Do something
            }
            return next(ctx, ...)
        }
    }
}
```

### State Middleware

```go
type StatefulMiddleware struct {
    counter int64
    mu      sync.Mutex
}

func (s *StatefulMiddleware) Middleware() PublisherMiddleware {
    return func(next PublisherFunc) PublisherFunc {
        return func(ctx context.Context, ...) error {
            s.mu.Lock()
            s.counter++
            count := s.counter
            s.mu.Unlock()
            
            logger.Info("message count", zap.Int64("count", count))
            return next(ctx, ...)
        }
    }
}
```

### Composable Middleware

```go
func Chain(middlewares ...PublisherMiddleware) PublisherMiddleware {
    return func(next PublisherFunc) PublisherFunc {
        for i := len(middlewares) - 1; i >= 0; i-- {
            next = middlewares[i](next)
        }
        return next
    }
}

// Use
combined := Chain(
    LoggingMiddleware(logger),
    MetricsMiddleware(),
    TracingMiddleware(tracer),
)
publisher.Use(combined)
```
