# NATS Package Production Readiness & Enhancements

This document outlines the features and improvements identified to make the `pkg/messaging/nats` package fully production-ready, robust, and scalable.

## 1. Pluggable Codecs (Flexibility)

**Current State**: 
The package hardcodes `json.Marshal` and `json.Unmarshal`. This limits the ability to use more efficient (Protobuf, MessagePack) or domain-specific serialization formats.

**Proposal**:
Introduce a `Codec` interface to abstract serialization.

### Detailed Implementation Sketch

#### 1. Codec Interface
Define a new interface in `pkg/messaging/nats/types.go`:
```go
// Codec defines the interface for message serialization and deserialization.
type Codec interface {
    // Encode marshals a value into bytes.
    Encode(v interface{}) ([]byte, error)
    // Decode unmarshals bytes into a value.
    Decode(data []byte, v interface{}) error
    // ContentType returns the content type string (e.g., "application/json").
    ContentType() string
}
```

#### 2. Default JSON Implementation
Provide a default implementation to maintain backward compatibility:
```go
type JSONCodec struct{}

func (c JSONCodec) Encode(v interface{}) ([]byte, error) {
    return json.Marshal(v)
}

func (c JSONCodec) Decode(data []byte, v interface{}) error {
    return json.Unmarshal(data, v)
}

func (c JSONCodec) ContentType() string {
    return "application/json"
}
```

#### 3. Client & Publisher Integration
- **Client Config**: Add `Codec` field to `Config` struct. In `NewNATSClient`, default to `&JSONCodec{}` if nil.
- **Publisher**: Update `Publish` methods to use `p.client.codec.Encode(data)` instead of `json.Marshal(data)`.
- **Subscriber**: The `HandlerFunc` currently receives `MessageEnvelope`. To support typed decoding, we can add a helper or update the handler signature in a v2. For now, the Subscriber usage remains similar, but the payload inside the envelope will be encoded by the Codec.

**Note**: To allow full flexibility (e.g., swapping JSON Envelopes for Protobuf), the Codec should apply to the **entire envelope payload**, meaning `msg.Data` (from NATS) is passed to `Decode`.

## 2. Publisher Resilience (Circuit Breaking & Retries)

**Current State**:
Publishing is direct. If the NATS server or JetStream is temporarily unavailable or slow, the application might block or fail immediately, potentially causing cascading failures in high-load scenarios.

**Proposal**:
Implement Circuit Breaking and Retry patterns to handle network instability and downstream pressure.

### Detailed Implementation Sketch

#### 1. Configuration
Add `ResilienceConfig` to the NATS `Config`:
```go
type ResilienceConfig struct {
    // Retry configuration
    MaxRetries     int           // e.g., 3
    InitialBackoff time.Duration // e.g., 100ms
    MaxBackoff     time.Duration // e.g., 1s
    
    // Circuit Breaker configuration (reference: sony/gobreaker)
    CBMaxRequests uint32        // Max requests in half-open state
    CBInterval    time.Duration // Cyclic period of closed state
    CBTimeout     time.Duration // Period of open state before half-open
}
```

#### 2. Circuit Breaker Integration
Integrate `github.com/sony/gobreaker` (or similar) into the `NATSPublisher`.

```go
import "github.com/sony/gobreaker"

type NATSPublisher struct {
    // ... existing fields ...
    breaker *gobreaker.CircuitBreaker
    resilience ResilienceConfig
}

// In NewPublisher:
settings := gobreaker.Settings{
    Name:        "NATSPublisher",
    Timeout:     config.CBTimeout,
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
        return counts.Requests >= 5 && failureRatio >= 0.6
    },
}
p.breaker = gobreaker.NewCircuitBreaker(settings)
```

#### 3. Publish Logic with Resilience
Wrap the core publish logic with both Retry and Circuit Breaker protections.

```go
func (p *NATSPublisher) Publish(...) error {
    // Define the operation to be retried
    operation := func() error {
        // Execute within Circuit Breaker
        _, err := p.breaker.Execute(func() (interface{}, error) {
            // Actual NATS call
            if opts.Async {
                return nil, p.client.Conn().Publish(subject, finalData)
            }
            return nil, p.client.Conn().Flush() // simplified
        })
        return err
    }

    return p.executeWithRetry(ctx, operation)
}

func (p *NATSPublisher) executeWithRetry(ctx context.Context, op func() error) error {
    var err error
    for i := 0; i <= p.resilience.MaxRetries; i++ {
        if err = op(); err == nil {
            return nil
        }
        
        // Fast fail if Circuit Breaker is open
        if err == gobreaker.ErrOpenState {
            return err
        }

        // Backoff
        delay := p.calculateBackoff(i)
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(delay):
            continue
        }
    }
    return fmt.Errorf("max retries exceeded: %w", err)
}
```

## 3. Metadata Injection (Observability & Routing)

**Current State**:
`MessageEnvelope` has a `Metadata` field, but `PublishOptions` does not support injecting arbitrary metadata (like `Tenant-ID`, `Correlation-ID`, or `Origin-Service`) at the call site.

**Proposal**:
Allow callers to inject custom metadata (e.g., Tenant ID, Correlation Source) per message.

### Detailed Implementation Sketch

#### 1. Update PublishOptions
Modify `pkg/messaging/nats/types.go` to include a Metadata map.

```go
type PublishOptions struct {
    Async    bool
    Timeout  time.Duration
    Metadata map[string]string // New field for custom headers/context
}
```

#### 2. Update Publisher Logic
In `pkg/messaging/nats/publisher.go`, merge the user-supplied metadata with the envelope's metadata.

```go
func (p *NATSPublisher) publish(ctx context.Context, ..., opts *PublishOptions) error {
    // ...
    envelope := MessageEnvelope{
        // ...
        Metadata: make(map[string]string),
    }

    // 1. Add User Metadata first
    if opts != nil && opts.Metadata != nil {
        for k, v := range opts.Metadata {
            envelope.Metadata[k] = v
        }
    }

    // 2. Add System/Middleware Metadata (e.g., Tracing)
    // OpenTelemetry injection often happens here or via middleware.
    // If using middleware, ensure it appends/overwrites as appropriate.
    otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(envelope.Metadata))

    // ... marshal and send
}
```

#### 3. Usage Example

```go
opts := &nats.PublishOptions{
    Metadata: map[string]string{
        "Tenant-ID": "tenant-123",
        "Source":    "payment-service",
    },
}
publisher.Publish(ctx, "orders.created", "OrderCreated", orderData, opts)
```

## 4. JetStream Reliability (Dead Letter Queues)

**Current State**:
`Subscriber` logs errors when processing JetStream messages fails. While JetStream handles redelivery, there is no explicit configuration or pattern for handling "poison messages" that repeatedly fail.

**Proposal**:
- **DLQ Configuration**: Allow configuring a Dead Letter Subject in `SubscribeOptions`.
- **Max Deliveries**: Expose JetStream's `MaxDeliver` option.
- **Handler Logic**: If processing fails `MaxDeliver` times, automatically publish the message to the DLQ subject (if NATS doesn't handle this via stream config) or ensure stream config is set up to route to a DLQ stream.

### Detailed Implementation Sketch

#### 1. Update SubscribeOptions
Add DLQ support to the subscription options.

```go
type SubscribeOptions struct {
    // ... existing fields ...
    
    // JetStream Specific
    MaxDeliver  int    // Max redelivery attempts before DLQ
    DLQSubject  string // Subject to publish to after max failures
}
```

#### 2. Enhanced Handler Logic
In `pkg/messaging/nats/subscriber.go`, update the `msgHandler` for JetStream.

```go
// Inside SubscribePush / SubscribePull
msgHandler := func(msg *nats.Msg) {
    // ... unmarshal and validation ...

    // Get current delivery count (JetStream metadata)
    meta, err := msg.Metadata()
    if err == nil {
        if opts.MaxDeliver > 0 && int(meta.NumDelivered) > opts.MaxDeliver {
            // Threshold exceeded - Send to DLQ
            if opts.DLQSubject != "" {
                s.publishToDLQ(ctx, opts.DLQSubject, msg)
                msg.Ack() // Ack original to remove from main stream
                return
            }
            // If no DLQ configured, let it Terminate (Ack + Log) or rely on Stream policy
            msg.Term() 
            return
        }
    }

    // Execute Handler
    err := h(ctx, msg.Subject, &envelope)
    if err != nil {
        // Nak to trigger redelivery if below threshold
        msg.Nak()
        return
    }
    
    msg.Ack()
}
```

#### 3. DLQ Publisher Helper
Add a helper to safely move the poison message.

```go
func (s *NATSSubscriber) publishToDLQ(ctx context.Context, subject string, originalMsg *nats.Msg) {
    // Wrap original message in a DLQ envelope with error context if possible
    // Or just forward the raw bytes
    
    dlqMsg := &nats.Msg{
        Subject: subject,
        Data:    originalMsg.Data,
        Header:  originalMsg.Header, // Preserve headers
    }
    // Add DLQ specific headers
    dlqMsg.Header.Add("X-Original-Subject", originalMsg.Subject)
    dlqMsg.Header.Add("X-Failure-Reason", "MaxDeliveredReached")

    s.client.Conn().PublishMsg(dlqMsg)
}
```

## 5. Strict Configuration Validation

**Current State**:
`Config` struct exists, but validation is minimal (mostly basic checks in `NewNATSClient` or connection failure at runtime).

**Proposal**:
- Implement a `Validate()` method on the `Config` struct.
- Check for valid URL formats.
- Ensure mutually exclusive options (e.g., Token vs Username/Password) are not both set.

### Detailed Implementation Sketch

#### 1. Validate Method
Add a `Validate` method to `Config` in `pkg/messaging/nats/client.go`.

```go
func (c Config) Validate() error {
    if c.URL == "" {
        return fmt.Errorf("NATS URL cannot be empty")
    }
    
    // Check mutual exclusion for auth
    if c.Token != "" && (c.Username != "" || c.Password != "") {
        return fmt.Errorf("cannot assume both Token and Username/Password authentication")
    }
    if c.CredsFile != "" && (c.Token != "" || c.Username != "") {
        return fmt.Errorf("cannot assume both Credentials File and other auth methods")
    }

    // Validate Timeouts
    if c.ConnectionTimeout < 0 {
        return fmt.Errorf("ConnectionTimeout cannot be negative")
    }

    return nil
}
```

#### 2. Integration
Call `Validate()` in `NewNATSClient`.

```go
func NewNATSClient(cfg Config, logger *zap.Logger) (*Client, error) {
    if logger == nil {
        return nil, fmt.Errorf("logger is required")
    }

    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("invalid configuration: %w", err)
    }

    // ... defaults ...
    return &Client{...}, nil
}
```

## 6. Testing & Simulation

**Current State**:
Tests rely on a live NATS server or mocks.

**Proposal**:
- **Embedded Server Tests**: Use `github.com/nats-io/nats-server/v2/test` to start real in-memory NATS servers for unit/integration tests to ensure stable CI environments without external dependencies.
- **Chaos Testing**: Add tests that simulate connection drops and slow consumers to verify resilience features.

## Completed Enhancements

- **Graceful Shutdown**: Implemented `sync.WaitGroup` in `Subscriber` to ensure active handlers complete before connection close.
