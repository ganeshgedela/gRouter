# NATS Event-Driven Microservice Template

This template demonstrates a **production-grade Event-Driven microservice** using `NATS JetStream` and the `gRouter` shared packages.

## Architecture

-   **Manager**: Handles lifecycle (Init, Start, Stop) and dependency injection
-   **NATS**: Uses `pkg/messaging/nats` for high-performance messaging
-   **Modular Services**: Demonstrates domain-driven design with separate service modules
-   **Config**: Centralized configuration management
-   **Logger**: Structured logging with Zap

## Directory Structure

```
.
├── cmd/
│   └── worker/             # Main entry point
├── internal/
│   ├── app/                # Application lifecycle management
│   │   └── app.go          # Main application orchestrator
│   ├── config/             # Service-specific configuration
│   │   └── config.go
│   └── pkg/                # Domain modules
│       ├── users/          # User management service
│       │   └── service.go
│       └── orders/         # Order processing service
│           └── service.go
├── Dockerfile              # Production Docker build
└── README.md               # This file
```

## Features

### 1. Modular Service Architecture
Each domain module (`users`, `orders`) is self-contained with:
- Event handlers
- Subscription management  
- Business logic

### 2. Production-Ready Patterns
- ✅ Queue Subscriptions for load balancing
- ✅ Durable consumers for reliability
- ✅ Graceful shutdown
- ✅ Health checks
- ✅ Structured logging
- ✅ Metrics (via middleware)

### 3. NATS Features
- Queue Groups for horizontal scaling
- JetStream for persistence
- Middleware chain (Recovery, Metrics, Logging, Tracing)

## How to Run

### Prerequisites
-   Go 1.22+
-   NATS Server with JetStream enabled
    ```bash
    nats-server -js
    ```

### Local Development

1.  Ensure `configs/config.yaml` has valid NATS settings:
    ```yaml
    nats:
      url: "nats://localhost:4222"
      middleware:
        metrics:
          enabled: true
        logging:
          enabled: true
    ```

2.  Run the worker:
    ```bash
    go run templates/nats-service/cmd/worker/main.go
    ```

3.  Publish test messages:
    ```bash
    # User created event
    nats pub users.created '{"user_id": "user-123", "email": "user@example.com", "name": "John Doe", "timestamp": "2024-01-08T10:00:00Z"}'
    
    # Order created event
    nats pub orders.created '{"order_id": "order-456", "user_id": "user-123", "amount": 99.99, "timestamp": "2024-01-08T10:00:00Z"}'
    ```

### Docker Build

```bash
docker build -f templates/nats-service/Dockerfile -t nats-worker:latest .
docker run -e NATS_URL=nats://host.docker.internal:4222 nats-worker:latest
```

## Adding New Services

To add a new domain service (e.g., `payments`):

1.  Create `internal/pkg/payments/service.go`:
    ```go
    package payments
    
    type Service struct {
        logger    *zap.Logger
        messenger *nats.Messenger
    }
    
    func (s *Service) Subscribe(ctx context.Context) error {
        // Register subscriptions
    }
    ```

2.  Register in `internal/app/app.go`:
    ```go
    a.paymentSvc = payments.NewService(a.deps.Logger, messenger)
    a.paymentSvc.Subscribe(ctx)
    ```

## Monitoring

- **Metrics**: Available at `/metrics` (if metrics server is enabled)
- **Logs**: Structured JSON logs to stdout
- **Health**: Implement via `HealthCheck()` method

## Production Checklist

- ✅ Graceful shutdown with context cancellation
- ✅ NATS reconnection handling (automatic)
- ✅ Dead letter queue support (configure in SubscribeOptions)
- ✅ Consumer acknowledgment policies
- ✅ Message deduplication (via JetStream)
- ✅ Horizontal scaling (queue groups)
