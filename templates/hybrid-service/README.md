# Hybrid Microservice Template

This template demonstrates a **production-grade Hybrid service** that bridges Synchronous (REST) and Asynchronous (NATS) communication.

## Use Case

"Backend for Frontend" (BFF) or Gateway services that:
- Accept user HTTP requests (synchronous)
- Dispatch NATS events for background processing (asynchronous)
- Enable event-driven architectures with REST API frontends

## Architecture

-   **Manager**: Orchestrates both `WebServer` and `NATS Messenger` lifecycle
-   **REST API**: Exposed via `pkg/web` for synchronous requests
-   **NATS Messaging**: Connected via `pkg/messaging/nats` for async events
-   **Config**: Centralized configuration management
-   **Logger**: Structured logging with Zap

## Directory Structure

```
.
├── cmd/
│   └── server/             # Main entry point
├── internal/
│   ├── app/                # Application lifecycle management
│   │   └── app.go          # Orchestrates REST + NATS
│   ├── config/             # Service-specific configuration
│   │   └── config.go
│   └── pkg/
│       └── api/            # HTTP handlers that publish events
│           └── handler.go
├── Dockerfile
└── README.md
```

## Features

### 1. Dual Protocol Support
- **REST API**: For synchronous client interactions
- **NATS Events**: For asynchronous processing and decoupling

### 2. Event-Driven Pattern
Example flow:
1. Client sends `POST /api/orders` (HTTP)
2. Handler creates order
3. Handler publishes `orders.created` event to NATS
4. Background workers process the event asynchronously

### 3. Production-Ready
- ✅ Graceful shutdown of both REST and NATS
- ✅ Health checks for both components
- ✅ Full middleware support (Metrics, Logging, Tracing)
- ✅ Connection health monitoring

## How to Run

### Prerequisites
-   Go 1.22+
-   NATS Server running

### 1. Start NATS Server
```bash
nats-server -js
```

### 2. Run the Hybrid Service
```bash
go run templates/hybrid-service/cmd/server/main.go
```

### 3. Test the API
```bash
# Health check
curl http://localhost:8080/api/health

# Create order (triggers NATS event)
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-123", "amount": 99.99}'
```

### 4. Monitor NATS Events
In another terminal, subscribe to the events:
```bash
nats sub "orders.>"
```

You should see the `orders.created` event published when you create an order via the REST API.

### Docker Build

```bash
docker build -f templates/hybrid-service/Dockerfile -t hybrid-service:latest .
docker run -e NATS_URL=nats://host.docker.internal:4222 -p 8080:8080 hybrid-service:latest
```

## API Endpoints

- `GET /api/health` - Health check
- `POST /api/orders` - Create order and publish event

## Architecture Flow

```
[Client] --HTTP POST--> [REST API Handler]
                              |
                              ├─> Store in DB (sync)
                              └─> Publish to NATS (async)
                                        |
                                        v
                              [Event Bus: orders.created]
                                        |
                                        v
                              [Background Workers]
                                  - Send email
                                  - Update inventory
                                  - Process payment
```

## Adding New Event-Driven Endpoints

1.  Create handler in `internal/pkg/api/handler.go`:
    ```go
    func (h *Handler) CreateUser(c *gin.Context) {
        // Business logic
        
        // Publish event
        h.messenger.Publish(ctx, "users.created", "user.created", userData, opts)
        
        // Return response
        c.JSON(201, response)
    }
    ```

2.  Register route in `RegisterRoutes()`:
    ```go
    api.POST("/users", h.CreateUser)
    ```

## Monitoring

- **REST Metrics**: Available at `/metrics`
- **NATS Health**: Checked in `/api/health` endpoint
- **Logs**: Structured JSON logs showing both HTTP and NATS events

## Production Checklist

- ✅ Graceful shutdown for both REST and NATS
- ✅ Connection retry and resilience
- ✅ Event acknowledgment handling
- ✅ Request timeout configuration
- ✅ Rate limiting on REST endpoints
- ✅ Circuit breaking for NATS publishing
- ✅ Health checks for both components
