# REST API Microservice Template

This template demonstrates a **production-grade REST API service** using `Gin` and the `pkg/web` shared package.

## Architecture

-   **Manager**: Handles lifecycle (Init, Start, Stop) and dependency injection
-   **Web Server**: Uses `pkg/web` for HTTP serving with full middleware support
-   **Modular Handlers**: Demonstrates domain-driven design with separate handler modules
-   **Config**: Centralized configuration management
-   **Logger**: Structured logging with Zap

## Directory Structure

```
.
├── cmd/
│   └── api/                # Main entry point
├── internal/
│   ├── app/                # Application lifecycle management
│   │   └── app.go          # Main application orchestrator
│   ├── config/             # Service-specific configuration
│   │   └── config.go
│   └── pkg/                # Domain modules
│       ├── users/          # User management handlers
│       │   └── handler.go
│       └── orders/         # Order management handlers
│           └── handler.go
├── Dockerfile              # Production Docker build
└── README.md               # This file
```

## Features

### 1. Modular Handler Architecture
Each domain module (`users`, `orders`) is self-contained with:
- RESTful route handlers
- CRUD operations
- Swagger documentation annotations

### 2. Production-Ready Patterns
- ✅ Automatic middleware (Recovery, Logging, Metrics, Tracing)
- ✅ Health endpoints (`/api/v1/health`)
- ✅ Graceful shutdown
- ✅ CORS support
- ✅ Rate limiting
- ✅ Request/Response compression

### 3. API Versioning
- Routes are versioned under `/api/v1`
- Easy to add new versions

## How to Run

### Prerequisites
-   Go 1.22+

### Local Development

1.  Ensure `configs/config.yaml` has valid Web settings:
    ```yaml
    web:
      port: 8080
      mode: debug
      middleware:
        metrics:
          enabled: true
        logging:
          enabled: true
    ```

2.  Run the API:
    ```bash
    go run templates/rest-service/cmd/api/main.go
    ```

3.  Test the endpoints:
    ```bash
    # Health check
    curl http://localhost:8080/api/v1/health
    
    # List users
    curl http://localhost:8080/api/v1/users
    
    # Get specific user
    curl http://localhost:8080/api/v1/users/1
    
    # List orders
    curl http://localhost:8080/api/v1/orders
    
    # Create order
    curl -X POST http://localhost:8080/api/v1/orders \
      -H "Content-Type: application/json" \
      -d '{"user_id": "1", "amount": 99.99}'
    ```

### Docker Build

```bash
docker build -f templates/rest-service/Dockerfile -t rest-api:latest .
docker run -p 8080:8080 rest-api:latest
```

## API Endpoints

### Users
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Orders
- `GET /api/v1/orders` - List all orders
- `GET /api/v1/orders/:id` - Get order by ID
- `POST /api/v1/orders` - Create new order
- `PUT /api/v1/orders/:id/status` - Update order status

## Adding New Handlers

To add a new domain handler (e.g., `products`):

1.  Create `internal/pkg/products/handler.go`:
    ```go
    package products
    
    type Handler struct {
        logger *zap.Logger
    }
    
    func (h *Handler) RegisterRoutes(router gin.IRouter) {
        // Register routes
    }
    ```

2.  Register in `internal/app/app.go`:
    ```go
    a.productHandler = products.NewHandler(a.deps.Logger)
    a.productHandler.RegisterRoutes(v1)
    ```

## Monitoring

- **Metrics**: Available at `/metrics` (Prometheus format)
- **Health**: Available at `/api/v1/health`
- **Logs**: Structured JSON logs to stdout

## Production Checklist

- ✅ Graceful shutdown with context cancellation
- ✅ Request timeout handling
- ✅ Rate limiting per client
- ✅ CORS configuration
- ✅ Security headers
- ✅ Request/Response logging
- ✅ Panic recovery
- ✅ OpenTelemetry tracing
