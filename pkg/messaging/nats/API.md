# NATS Package API Reference

## Core Interfaces

### Messenger

The `Messenger` is the primary interface for accessing NATS functionality.

```go
type Messenger struct {
    Client     *Client
    Publisher  Publisher
    Subscriber Subscriber
}
```

#### Methods

**`InitFromConfig(cfg *config.Config, logger *zap.Logger) error`**
- Initializes messenger from global configuration
- Validates configuration (URL, timeouts, middleware settings)
- Connects to NATS server
- Returns error if configuration is invalid or connection fails

**`Start() error`**
- Starts the messenger and initializes pub/sub with middleware chains
- Must be called after `InitFromConfig`
- Registers all configured middleware

**`Stop(ctx context.Context) error`**
- Gracefully stops the messenger
- Drains all subscriptions
- Closes NATS connection

---

## Publisher Interface

```go
type Publisher interface {
    // Synchronous publish
    Publish(ctx context.Context, subject, msgType string, data interface{}, opts *PublishOptions) error
    
    // Request-Reply pattern
    Request(ctx context.Context, subject, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error)
    
    // JetStream publish (synchronous)
    PublishJS(ctx context.Context, subject, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error)
    
    // JetStream publish (asynchronous)
    PublishAsyncJS(ctx context.Context, subject, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error)
    
    // Middleware registration
    Use(middleware ...PublisherMiddleware)
    UseRequest(middleware ...RequestMiddleware)
    UseJS(middleware ...JSPublisherMiddleware)
    UseAsyncJS(middleware ...JSAsyncPublisherMiddleware)
}
```

### Publish Options

```go
type PublishOptions struct {
    Headers  map[string]string  // NATS headers for routing/filtering
    Metadata map[string]string  // Envelope metadata (tracing, etc.)
}
```

### Usage Examples

**Simple Publish:**
```go
err := publisher.Publish(ctx, "orders.created", "OrderCreated", orderData, nil)
```

**Publish with Headers:**
```go
opts := &PublishOptions{
    Headers: map[string]string{
        "priority": "high",
    },
    Metadata: map[string]string{
        "correlation_id": "12345",
    },
}
err := publisher.Publish(ctx, "orders.created", "OrderCreated", orderData, opts)
```

**Request-Reply:**
```go
response, err := publisher.Request(ctx, "inventory.check", "CheckInventory", 
    checkRequest, 5*time.Second, nil)
if err != nil {
    // Handle timeout or error
}
// Process response.Data
```

**JetStream (At-Least-Once):**
```go
ack, err := publisher.PublishJS(ctx, "orders.persistent", "Order", orderData, nil)
if err != nil {
    // Handle publish failure
}
log.Printf("Published to stream, seq: %d", ack.Sequence)
```

---

## Subscriber Interface

```go
type Subscriber interface {
    // Core NATS subscription
    Subscribe(subject string, handler HandlerFunc, opts *SubscribeOptions) (*nats.Subscription, error)
    
    // Queue group subscription (load balancing)
    SubscribeQueue(subject, queue string, handler HandlerFunc, opts *SubscribeOptions) (*nats.Subscription, error)
    
    // JetStream push consumer
    SubscribePush(subject string, handler HandlerFunc, opts ...nats.SubOpt) (*nats.Subscription, error)
    
    // JetStream pull consumer (worker pattern)
    SubscribePull(subject string, handler HandlerFunc, opts ...nats.SubOpt) (*nats.Subscription, error)
    
    // Middleware registration
    Use(middleware ...SubscriberMiddleware)
    
    // Graceful shutdown
    Close(ctx context.Context) error
}
```

### Subscribe Options

```go
type SubscribeOptions struct {
    QueueGroup string  // Queue group name for load balancing
    MaxWorkers int     // Concurrency limit (default: unlimited)
}
```

### Handler Function

```go
type HandlerFunc func(ctx context.Context, subject string, envelope *MessageEnvelope) error
```

### Usage Examples

**Simple Subscribe:**
```go
sub, err := subscriber.Subscribe("orders.created", func(ctx context.Context, subject string, env *MessageEnvelope) error {
    var order Order
    if err := json.Unmarshal(env.Data, &order); err != nil {
        return err
    }
    // Process order
    return nil
}, nil)
```

**Queue Group (Load Balancing):**
```go
opts := &SubscribeOptions{
    QueueGroup: "order-processors",
    MaxWorkers: 10,  // Process max 10 messages concurrently
}
sub, err := subscriber.Subscribe("orders.created", handler, opts)
```

**JetStream Push (Durable):**
```go
sub, err := subscriber.SubscribePush("orders.critical", handler, 
    nats.Durable("order-processor"),
    nats.ManualAck(),
)
```

**JetStream Pull (Worker Pattern):**
```go
sub, err := subscriber.SubscribePull("batch.jobs", handler,
    nats.Durable("job-worker"),
    nats.PullMaxWaiting(128),
)
```

---

## Client Interface

```go
type Client struct {
    conn   *nats.Conn
    js     nats.JetStreamContext
    logger *zap.Logger
    source string
    config Config
}
```

### Methods

**`Connect() error`**
- Establishes connection to NATS server
- Applies TLS/authentication configuration
- Sets up reconnection handlers

**`Close() error`**
- Drains and closes the connection
- Safe to call multiple times

**`IsConnected() bool`**
- Returns true if connection is active

**`Conn() *nats.Conn`**
- Returns underlying NATS connection for advanced usage

**`JetStream() nats.JetStreamContext`**
- Returns JetStream context for stream management

---

## Configuration

### Complete Config Structure

```go
type Config struct {
    // Connection
    URL               string        // NATS server URL
    MaxReconnects     int           // Max reconnection attempts
    ReconnectWait     time.Duration // Delay between reconnects
    ConnectionTimeout time.Duration // Initial connection timeout
    DrainTimeout      time.Duration // Graceful shutdown timeout
    
    // Authentication
    Token     string // Token-based auth
    Username  string // User/password auth
    Password  string
    CredsFile string // NATS 2.0 JWT/NKey credentials
    
    // TLS
    UseTLS     bool   // Enable TLS
    SkipVerify bool   // Skip certificate verification (dev only)
    CAFile     string // CA certificate
    CertFile   string // Client certificate
    KeyFile    string // Client key
    
    // Middleware
    Middleware NATSMiddlewareConfig
}
```

### Middleware Configuration

```go
type NATSMiddlewareConfig struct {
    Recovery       MiddlewareState
    Metrics        MetricsConfig
    Tracing        MiddlewareState
    Logging        LoggingConfig
    CircuitBreaker CircuitBreakerConfig
    Retry          RetryConfig
    Timeout        TimeoutConfig
    RateLimit      RateLimitConfig
}

type CircuitBreakerConfig struct {
    Enabled       bool
    MaxRequests   uint32        // Max requests in half-open state
    Interval      time.Duration // Interval for counting failures
    Timeout       time.Duration // Timeout before trying half-open
    TripThreshold uint32        // Consecutive failures to trip
}

type RetryConfig struct {
    Enabled         bool
    MaxAttempts     int
    InitialInterval time.Duration
    Multiplier      float64       // Backoff multiplier
    MaxInterval     time.Duration // Max backoff duration
}

type RateLimitConfig struct {
    Enabled           bool
    RequestsPerSecond int
    Burst             int  // Token bucket burst size
}
```

---

## Message Envelope

All messages are wrapped in a standardized envelope:

```go
type MessageEnvelope struct {
    ID        string            `json:"id"`        // Unique message ID (UUID)
    Type      string            `json:"type"`      // Message type/schema identifier
    Timestamp time.Time         `json:"timestamp"` // Message creation time
    Source    string            `json:"source"`    // Originating service
    Reply     string            `json:"reply"`     // Reply-to subject (for request-reply)
    Data      json.RawMessage   `json:"data"`      // Actual payload
    Metadata  map[string]string `json:"metadata"`  // Contextual metadata (trace IDs, etc.)
}
```

### Envelope Helpers

**Creating an Envelope:**
```go
env := &MessageEnvelope{
    ID:        uuid.New().String(),
    Type:      "OrderCreated",
    Timestamp: time.Now(),
    Source:    "order-service",
    Data:      jsonData,
}
```

**Extracting Data:**
```go
var order Order
if err := json.Unmarshal(env.Data, &order); err != nil {
    return err
}
```

---

## Error Handling

### Common Errors

- `ErrNotConnectedE`: Client is not connected to NATS
- `ErrTimeout`: Operation timed out
- `ErrNoResponders`: No service listening on request subject
- `service temporarily unavailable`: Circuit breaker is open

### Error Handling Best Practices

```go
err := publisher.Publish(ctx, subject, msgType, data, nil)
if err != nil {
    if err == nats.ErrNoResponders {
        // No subscriber, queue for retry or log
    } else if err.Error() == "service temporarily unavailable" {
        // Circuit breaker open, back off
    } else {
        // Other error, handle appropriately
    }
}
```

---

## Thread Safety

- **Client**: Thread-safe for concurrent publish/subscribe operations
- **Publisher**: Thread-safe, can be used from multiple goroutines
- **Subscriber**: Thread-safe subscription management
- **Validator**: Thread-safe schema registration with `sync.RWMutex`

---

## Best Practices

### Connection Management
```go
// Use a single client per application
client, _ := NewNATSClient(cfg, logger)
defer client.Close()

// Reuse publisher/subscriber instances
pub := NewPublisher(client)
sub := NewSubscriber(client)
```

### Context Usage
```go
// Always pass context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := publisher.Publish(ctx, subject, msgType, data, nil)
```

### Graceful Shutdown
```go
// In your service shutdown handler
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := messenger.Stop(ctx); err != nil {
    log.Printf("Shutdown error: %v", err)
}
```

### Middleware Ordering
```go
// Register middleware in order of execution
// Recovery should be outermost (first)
publisher.Use(
    RecoveryMiddleware(logger),      // Catches panics
    MetricsMiddleware(),              // Records metrics
    LoggingMiddleware(logger),        // Logs operations
    CircuitBreakerMiddleware(cb),     // Fails fast
    RetryMiddleware(retryCfg),        // Retries failures
)
```
