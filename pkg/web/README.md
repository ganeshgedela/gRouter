# Web Framework Package (`pkg/web`)

The `pkg/web` package provides a production-ready, generic web server framework based on [Gin](https://github.com/gin-gonic/gin). It is designed to be modular, secure, and observable, allowing services to easily register routes and handlers.

## Features

### 1. Observability
- **Prometheus Metrics**: Automatically exposes standard HTTP metrics (request count, latency, size) at `/metrics`.
- **OpenTelemetry Tracing**: Integrated tracing middleware to propagate trace contexts.
- **Health Checks**: Built-in `HealthManager` exposing `/health/live` and `/health/ready` endpoints.

### 2. Security
- **TLS Support**: Configurable TLS termination with certificate and key files.
- **CORS**: Flexible Cross-Origin Resource Sharing configuration.
- **Security Headers**: Automatic injection of security headers (HSTS, X-Frame-Options, XSS-Protection, CSP) via `secure` middleware.

### 3. Resilience
- **Rate Limiting**: IP-based rate limiting with configurable requests per second and burst capacity.
- **Request ID**: Automatically assigns and logs a unique `X-Request-ID` for every request.
- **Graceful Shutdown**: Handles OS signals to shut down the server gracefully, waiting for active connections to complete.

## Usage

### 1. Initialization
Initialize the server with a configuration and a logger.

```go
import (
    "github.com/ganesh/grouter/pkg/web"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()
    cfg := web.DefaultConfig()
    cfg.Port = 8080

    server := web.NewServer(cfg, logger)
    
    // ... register services ...

    if err := server.Start(); err != nil {
        logger.Fatal("Failed to start server", zap.Error(err))
    }
}
```

### 2. Registering Services
Services must implement the `web.Service` interface:

```go
type Service interface {
    RegisterRoutes(router *gin.RouterGroup)
}
```

Example implementation:

```go
type MyService struct {
    log *zap.Logger
}

func (s *MyService) RegisterRoutes(router *gin.RouterGroup) {
    router.GET("/hello", s.HelloHandler)
}

func (s *MyService) HelloHandler(c *gin.Context) {
    c.JSON(200, gin.H{"message": "Hello World"})
}
```

Register the service with the server:

```go
myService := &MyService{log: logger}
server.RegisterService(myService)
```

### 3. Adding Health Checks
You can register custom health checks for your services.

```go
server.AddReadinessCheck("database", func() error {
    return db.Ping()
})

server.AddLivenessCheck("worker", func() error {
    if !worker.IsRunning() {
        return fmt.Errorf("worker stopped")
    }
    return nil
})
```

## Configuration

The server is fully configurable via the `WebConfig` struct or YAML configuration.

```yaml
web:
  port: 8080
  read_timeout: 10s
  write_timeout: 10s
  shutdown_timeout: 5s
  mode: "release" # debug, release, test
  
  metrics:
    enabled: true
    path: "/metrics"
    
  tracing:
    enabled: true
    service_name: "grouter-web"
    
  tls:
    enabled: true
    cert_file: "/path/to/cert.pem"
    key_file: "/path/to/key.pem"
    
  cors:
    enabled: true
    allowed_origins: ["https://example.com"]
    allowed_methods: ["GET", "POST"]
    allowed_headers: ["Origin", "Content-Type"]
    allow_credentials: true
    max_age: 3600
    
  security:
    enabled: true
    xss_protection: "1; mode=block"
    content_type_nosniff: "nosniff"
    x_frame_options: "DENY"
    hsts_max_age: 31536000
    hsts_exclude_subdomains: true
    content_security_policy: "default-src 'self'"
    referrer_policy: "strict-origin-when-cross-origin"
    
  rate_limit:
    enabled: true
    requests_per_second: 100
    burst: 200
```

## Build & Test

### Bazel Support

This project uses Bazel for building and testing.

#### Prerequisites
- [Bazel](https://bazel.build/install) installed (recommended: use `bazelisk`).
- Go SDK installed.

#### Build Steps

1.  **Update Dependencies**:
    If you add new dependencies to `go.mod`, update the Bazel repositories:
    ```bash
    bazel run //:gazelle-update-repos
    ```
    Then update the `BUILD.bazel` files:
    ```bash
    bazel run //:gazelle
    ```

2.  **Build the Web Package**:
    ```bash
    bazel build //pkg/web
    ```
    To build the natsdemosvc service:
    ```bash
    bazel build //pkg/services/natsdemosvc
    ```

3.  **Run Tests**:
    To run the unit tests for the web package:
    ```bash
    bazel test //pkg/web:web_test
    ```
    To run tests for the natsdemosvc service:
    ```bash
    bazel test //pkg/services/natsdemosvc:natsdemosvc_test
    ```
    To run all tests in the workspace:
    ```bash
    bazel test //...
    ```

## Architecture

- **Server**: The core struct that wraps `gin.Engine` and `http.Server`. It manages the lifecycle and middleware chain.
- **Middleware**:
    - `LoggerMiddleware`: Logs request details and status.
    - `Recovery`: Recovers from panics.
    - `MetricsMiddleware`: Prometheus instrumentation.
    - `otelgin`: OpenTelemetry tracing.
    - `cors`: Handles CORS preflight and headers.
    - `secure`: Adds security headers.
    - `RateLimitMiddleware`: Token bucket rate limiting per IP.
    - `RequestIDMiddleware`: Injects unique request IDs.
- **HealthManager**: A thread-safe manager for liveness and readiness probes.

## Future Roadmap

The following features are planned for future releases to enhance the framework's capabilities:

1.  **Authentication & Authorization**
    - Implement middleware for JWT and OAuth2 validation.
    - Add role-based access control (RBAC) support.

2.  **API Documentation**
    - Integrate `swaggo/gin-swagger` to automatically generate OpenAPI/Swagger documentation from code comments.

3.  **Distributed Rate Limiting**
    - Replace the in-memory rate limiter with a Redis-backed implementation to support distributed deployments (multiple replicas of gRouter).

4.  **Response Caching**
    - Add middleware for caching responses (in-memory or Redis) to improve performance for static or expensive data.

5.  **WebSocket Support**
    - Add native support for WebSocket upgrades and connection management for real-time features.

6.  **Dynamic Configuration**
    - Support hot-reloading of configuration (e.g., log levels, rate limits) without restarting the server.

