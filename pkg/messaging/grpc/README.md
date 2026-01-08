# gRPC Package - Quick Start Guide

## gRPC Package

Production-ready gRPC server and client with comprehensive middleware support, TLS/mTLS, and authentication.

## Features

### Core
- ✅ gRPC server with reflection
- ✅ gRPC client with connection management
- ✅ gRPC-Gateway support (HTTP/JSON to gRPC)
- ✅ Middleware/interceptor chain

### Security
- ✅ **TLS/mTLS** - Encrypted communication
- ✅ **JWT Authentication** - Token-based auth
- ✅ **API Key Authentication** - Service-to-service auth
- ✅ Minimum TLS 1.2 enforcement

### Middleware (7 Types)
1. **Recovery** - Panic recovery
2. **Logging** - Structured logging
3. **Metrics** - Prometheus integration
4. **Tracing** - OpenTelemetry support
5. **Retry** - Exponential backoff (client)
6. **Timeout** - Request timeouts
7. **Rate Limit + Circuit Breaker** - Traffic control

### Production Ready
- ✅ 77%+ test coverage
- ✅ Comprehensive documentation
- ✅ Example implementations
- ✅ Best practices built-in

## Quick Start

### 1. Server

```go
package main

import (
    "grouter/pkg/messaging/grpc"
    "grouter/pkg/messaging/grpc/interceptors"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()
    
    // Create server with interceptors
    server := grpc.NewServer(logger,
        grpc.WithPort(9090),
        grpc.WithReflection(),
        grpc.WithUnaryInterceptor(interceptors.Recovery(logger)),
        grpc.WithUnaryInterceptor(interceptors.Logging(logger)),
        grpc.WithUnaryInterceptor(interceptors.Metrics()),
    )
    
    // Register your services
    // pb.RegisterGreeterServer(server, &myService{})
    
    // Start server
    server.Start()
}
```

### 2. Client

```go
package main

import (
    "grouter/pkg/messaging/grpc"
    "grouter/pkg/messaging/grpc/interceptors"
    "time"
)

func main() {
    config := grpc.ClientConfig{
        Target:  "localhost:9090",
        Timeout: 10 * time.Second,
    }
    
    client, _ := grpc.NewClient(config, logger,
        grpc.WithClientUnaryInterceptor(interceptors.LoggingClient(logger)),
        grpc.WithClientUnaryInterceptor(interceptors.Retry(interceptors.DefaultRetryConfig())),
    )
    defer client.Close()
    
    // Use connection
    // pb.NewGreeterClient(client.GetConn())
}

## Security

### TLS/mTLS

```go
// Server with TLS
server := grpc.NewServer(logger,
    grpc.WithPort(8443),
    grpc.WithTLS("/certs/server.pem", "/certs/server-key.pem"),
)

// Server with mutual TLS (client cert required)
server := grpc.NewServer(logger,
    grpc.WithPort(8443),
    grpc.WithMTLS("/certs/server.pem", "/certs/server-key.pem", "/certs/ca.pem"),
)

// Client with TLS
client, err := grpc.NewClient(config, logger,
    grpc.WithClientTLS("api.example.com", "/certs/ca.pem"),
)

// Client with mTLS
client, err := grpc.NewClient(config, logger,
    grpc.WithClientMTLS("api.example.com", 
        "/certs/client.pem", "/certs/client-key.pem", "/certs/ca.pem"),
)
```

### JWT Authentication

```go
import "grouter/pkg/messaging/grpc/interceptors"

// Create JWT authenticator
auth := interceptors.NewJWTAuthenticator("your-secret-key", "your-issuer")

// Add to server
server := grpc.NewServer(logger,
    grpc.WithUnaryInterceptor(interceptors.AuthUnaryInterceptor(
        interceptors.AuthConfig{
            Enabled: true,
            Authenticator: auth,
            SkipMethods: []string{"/health", "/metrics"}, // Skip auth for these
        },
    )),
)

// Client usage
token := "eyJhbGciOiJIUzI1NiIs..." // Your JWT token
md := metadata.New(map[string]string{
    "authorization": "Bearer " + token,
})
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

### API Key Authentication

```go
// Create API key authenticator
apiAuth := interceptors.NewAPIKeyAuthenticator(map[string]interceptors.APIKeyInfo{
    "service-key-123": {
        UserID: "service-a",
        Email:  "service-a@example.com",
        Roles:  []string{"service"},
    },
})

// Add to server
server := grpc.NewServer(logger,
    grpc.WithUnaryInterceptor(interceptors.AuthUnaryInterceptor(
        interceptors.AuthConfig{
            Enabled: true,
            Authenticator: apiAuth,
        },
    )),
)

// Client usage
md := metadata.New(map[string]string{
    "x-api-key": "service-key-123",
})
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

### Access User Info in Handlers

```go
func (s *Server) MyMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    // Get authenticated user info
    userID, ok := interceptors.GetUserID(ctx)
    if !ok {
        return nil, status.Error(codes.Unauthenticated, "no user ID")
    }
    
    email, _ := interceptors.GetUserEmail(ctx)
    roles, _ := interceptors.GetUserRoles(ctx)
    
    log.Printf("Request from user: %s (%s), roles: %v", userID, email, roles)
    
    // Your logic here
    return &pb.Response{}, nil
}
```

## Interceptors

| Interceptor | Type | Purpose |
|-------------|------|---------|
| Recovery | Server | Catch panics |
| Logging | Both | Request logging |
| Metrics | Both | Prometheus metrics |
| Tracing | Both | OpenTelemetry tracing |
| Retry | Client | Exponential backoff |
| Timeout | Both | Deadline enforcement |
| RateLimit | Server | Throttling |
| CircuitBreaker | Client | Fault tolerance |

## Configuration

### Server Options
- `WithPort(int)` - Set port
- `WithReflection()` - Enable reflection
- `WithUnaryInterceptor(...)` - Add interceptor
- `WithGateway(ctx, port)` - Enable REST gateway

### Client Options
- `WithClientUnaryInterceptor(...)` - Add interceptor

## Examples

See `/pkg/messaging/grpc/example/main.go` for complete examples.

## Documentation

- **ARCHITECTURE.md** - Design & architecture
- **grpc_learning.md** - Detailed usage guide
- **interceptors/** - Interceptor implementations

## Testing

```bash
go test ./pkg/messaging/grpc/...
```

## Metrics

Server metrics:
```
grpc_server_requests_total{method, code}
grpc_server_request_duration_seconds{method, code}
```

Client metrics:
```
grpc_client_requests_total{method, code}
grpc_client_request_duration_seconds{method, code}
```

Access at `/metrics` endpoint (requires Prometheus handler).
