# gRPC Package - Complete Learning Guide

## Table of Contents
1. [Introduction](#introduction)
2. [Basic Server Setup](#basic-server-setup)
3. [Basic Client Setup](#basic-client-setup)
4. [Interceptor Usage](#interceptor-usage)
5. [Production Patterns](#production-patterns)
6. [Integration Examples](#integration-examples)

---

## Introduction

The gRPC package provides production-ready server and client wrappers with comprehensive middleware (interceptor) support. This guide shows practical usage patterns.

---

## Basic Server Setup

### Minimal Server

```go
package main

import (
    "grouter/pkg/messaging/grpc"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()
    
    server := grpc.NewServer(logger,
        grpc.WithPort(9090),
    )
    
    // Register services
    // pb.RegisterYourServiceServer(server, &yourImpl{})
    
    server.Start()
}
```

### Server with Reflection (for grpcurl)

```go
server := grpc.NewServer(logger,
    grpc.WithPort(9090),
    grpc.WithReflection(), // Enable reflection
)
```

Test with grpcurl:
```bash
grpcurl -plaintext localhost:9090 list
```

### Server with Gateway (REST ↔ gRPC)

```go
import "context"

ctx := context.Background()

server := grpc.NewServer(logger,
    grpc.WithPort(9090),
    grpc.WithGateway(ctx, 8080), // REST on 8080, gRPC on 9090
)

// Register gateway handlers
server.RegisterGatewayHandler(pb.RegisterGreeterHandlerFromEndpoint)
```

Access via REST:
```bash
curl http://localhost:8080/v1/greeter/hello
```

---

## Basic Client Setup

### Simple Client

```go
import (
    "grouter/pkg/messaging/grpc"
    "time"
)

func NewGreeterClient(logger *zap.Logger) (pb.GreeterClient, error) {
    config := grpc.ClientConfig{
        Target:  "localhost:9090",
        Timeout: 10 * time.Second,
    }
    
    client, err := grpc.NewClient(config, logger)
    if err != nil {
        return nil, err
    }
    
    return pb.NewGreeterClient(client.GetConn()), nil
}
```

### Client with Keepalive

```go
config := grpc.ClientConfig{
    Target:          "localhost:9090",
    Timeout:         10 * time.Second,
    KeepAliveTime:   30 * time.Second,
    KeepAliveTimeout: 10 * time.Second,
}

client, _ := grpc.NewClient(config, logger)
```

---

## Interceptor Usage

### Server Interceptors

#### 1. Recovery (Panic Protection)

**Must be first!**

```go
import "grouter/pkg/messaging/grpc/interceptors"

server := grpc.NewServer(logger,
    grpc.WithUnaryInterceptor(interceptors.Recovery(logger)),
    grpc.WithStreamInterceptor(interceptors.RecoveryStream(logger)),
)
```

**Example panic handling:**
```go
func (s *Service) RiskyMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    panic("unexpected error") // Caught by Recovery interceptor
    // Returns: status.Error(codes.Internal, "internal server error: unexpected error")
}
```

#### 2. Logging

```go
server := grpc.NewServer(logger,
    grpc.WithUnaryInterceptor(interceptors.Logging(logger)),
    grpc.WithStreamInterceptor(interceptors.LoggingStream(logger)),
)
```

**Log output:**
```json
{
  "level": "debug",
  "msg": "gRPC request received",
  "method": "/greeter.Greeter/SayHello",
  "peer": "127.0.0.1:54321"
}
{
  "level": "debug",
  "msg": "gRPC request completed",
  "method": "/greeter.Greeter/SayHello",
  "duration": "45ms",
  "code": "OK"
}
```

#### 3. Metrics

```go
server := grpc.NewServer(logger,
    grpc.WithUnaryInterceptor(interceptors.Metrics()),
    grpc.WithStreamInterceptor(interceptors.MetricsStream()),
)
```

**Expose metrics:**
```go
import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

http.Handle("/metrics", promhttp.Handler())
go http.ListenAndServe(":2112", nil)
```

**Query metrics:**
```bash
curl http://localhost:2112/metrics | grep grpc_server
```

#### 4. Tracing

```go
import "go.opentelemetry.io/otel"

tracer := otel.Tracer("my-service")

server := grpc.NewServer(logger,
    grpc.WithUnaryInterceptor(interceptors.Tracing(tracer)),
    grpc.WithStreamInterceptor(interceptors.TracingStream(tracer)),
)
```

#### 5. Rate Limiting

```go
// Global rate limit: 100 req/s, burst of 10
server := grpc.NewServer(logger,
    grpc.WithUnaryInterceptor(interceptors.RateLimiter(100, 10)),
)
```

**Per-method rate limits:**
```go
limits := interceptors.NewMethodRateLimits()
limits.Add("/greeter.Greeter/SayHello", 50, 5)      // 50 req/s
limits.Add("/greeter.Greeter/SayHelloAgain", 100, 10) // 100 req/s

server := grpc.NewServer(logger,
    grpc.WithUnaryInterceptor(
        interceptors.PerMethodRateLimiterInterceptor(limits, 10, 5), // default 10 req/s
    ),
)
```

#### 6. Timeout

```go
import "time"

// Enforce 30 second timeout on all requests
server := grpc.NewServer(logger,
    grpc.WithUnaryInterceptor(interceptors.Timeout(30 * time.Second)),
)
```

### Client Interceptors

#### 1. Retry

```go
import "grouter/pkg/messaging/grpc/interceptors"

// Use default config
client, _ := grpc.NewClient(config, logger,
    grpc.WithClientUnaryInterceptor(
        interceptors.Retry(interceptors.DefaultRetryConfig()),
    ),
)
```

**Custom retry config:**
```go
retryConfig := interceptors.RetryConfig{
    MaxAttempts: 5,
    InitialWait: 100 * time.Millisecond,
    MaxWait:     5 * time.Second,
    Multiplier:  2.0,
    RetryableCodes: map[codes.Code]bool{
        codes.Unavailable:       true,
        codes.ResourceExhausted: true,
        codes.Aborted:          true,
        codes.DeadlineExceeded: true,
    },
}

client, _ := grpc.NewClient(config, logger,
    grpc.WithClientUnaryInterceptor(interceptors.Retry(retryConfig)),
)
```

**Backoff sequence:**
```
Attempt 1: 0ms
Attempt 2: 100ms
Attempt 3: 200ms
Attempt 4: 400ms
Attempt 5: 800ms
```

#### 2. Circuit Breaker

```go
// Create circuit breaker: 5 failures → open, 60s timeout, 1 test request
cb := interceptors.NewCircuitBreaker(5, 60*time.Second, 1)

client, _ := grpc.NewClient(config, logger,
    grpc.WithClientUnaryInterceptor(
        interceptors.CircuitBreakerInterceptor(cb),
    ),
)
```

**Circuit breaker states:**
- **Closed:** Normal operation, requests allowed
- **Open:** Too many failures, all requests rejected (returns error immediately)
- **Half-Open:** Testing recovery, limited requests allowed

#### 3. Client Timeout

```go
client, _ := grpc.NewClient(config, logger,
    grpc.WithClientUnaryInterceptor(
        interceptors.TimeoutClient(10 * time.Second),
    ),
)
```

#### 4. Client Logging & Metrics

```go
client, _ := grpc.NewClient(config, logger,
    grpc.WithClientUnaryInterceptor(interceptors.LoggingClient(logger)),
    grpc.WithClientUnaryInterceptor(interceptors.MetricsClient()),
)
```

---

## Production Patterns

### Full Server Stack

```go
import (
    "grouter/pkg/messaging/grpc"
    "grouter/pkg/messaging/grpc/interceptors"
    "go.opentelemetry.io/otel"
    "go.uber.org/zap"
    "time"
)

func NewProductionServer(logger *zap.Logger) *grpc.Server {
    tracer := otel.Tracer("my-service")
    
    return grpc.NewServer(logger,
        grpc.WithPort(9090),
        grpc.WithReflection(),
        
        // Unary interceptors (order matters!)
        grpc.WithUnaryInterceptor(interceptors.Recovery(logger)),          // 1. Catch panics
        grpc.WithUnaryInterceptor(interceptors.Logging(logger)),           // 2. Log requests
        grpc.WithUnaryInterceptor(interceptors.Metrics()),                 // 3. Record metrics
        grpc.WithUnaryInterceptor(interceptors.Tracing(tracer)),           // 4. Distributed tracing
        grpc.WithUnaryInterceptor(interceptors.RateLimiter(100, 10)),      // 5. Rate limiting
        grpc.WithUnaryInterceptor(interceptors.Timeout(30*time.Second)),   // 6. Timeout
        
        // Stream interceptors
        grpc.WithStreamInterceptor(interceptors.RecoveryStream(logger)),
        grpc.WithStreamInterceptor(interceptors.LoggingStream(logger)),
        grpc.WithStreamInterceptor(interceptors.MetricsStream()),
        grpc.WithStreamInterceptor(interceptors.TracingStream(tracer)),
    )
}
```

### Full Client Stack

```go
func NewProductionClient(target string, logger *zap.Logger) (*grpc.Client, error) {
    tracer := otel.Tracer("my-client")
    cb := interceptors.NewCircuitBreaker(5, 60*time.Second, 1)
    
    config := grpc.ClientConfig{
        Target:           target,
        Timeout:          10 * time.Second,
        KeepAliveTime:    30 * time.Second,
        KeepAliveTimeout: 10 * time.Second,
    }
    
    return grpc.NewClient(config, logger,
        // Client interceptors (order matters!)
        grpc.WithClientUnaryInterceptor(interceptors.LoggingClient(logger)),                 // 1. Log
        grpc.WithClientUnaryInterceptor(interceptors.MetricsClient()),                       // 2. Metrics
        grpc.WithClientUnaryInterceptor(interceptors.TracingClient(tracer)),                 // 3. Tracing
        grpc.WithClientUnaryInterceptor(interceptors.CircuitBreakerInterceptor(cb)),         // 4. Circuit breaker
        grpc.WithClientUnaryInterceptor(interceptors.Retry(interceptors.DefaultRetryConfig())), // 5. Retry
        grpc.WithClientUnaryInterceptor(interceptors.TimeoutClient(10*time.Second)),         // 6. Timeout
    )
}
```

### Graceful Shutdown

```go
type App struct {
    server *grpc.Server
}

func (a *App) Start() {
    a.server.Start()
}

func (a *App) Stop() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    a.server.Stop(ctx) // Gracefully drain in-flight requests
}

func main() {
    app := &App{server: NewProductionServer(logger)}
    
    go app.Start()
    
    // Wait for signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    
    app.Stop()
}
```

---

## Integration Examples

### Complete Greeter Service

**1. Define proto (helloworld.proto):**
```protobuf
syntax = "proto3";

package helloworld;

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

**2. Server implementation:**
```go
package main

import (
    "context"
    "grouter/pkg/messaging/grpc"
    "grouter/pkg/messaging/grpc/interceptors"
    pb "grouter/api/proto/helloworld"
    "go.uber.org/zap"
)

type greeterServer struct {
    pb.UnimplementedGreeterServer
    logger *zap.Logger
}

func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
    s.logger.Info("Received request", zap.String("name", req.Name))
    return &pb.HelloReply{Message: "Hello " + req.Name}, nil
}

func main() {
    logger, _ := zap.NewProduction()
    
    server := grpc.NewServer(logger,
        grpc.WithPort(9090),
        grpc.WithReflection(),
        grpc.WithUnaryInterceptor(interceptors.Recovery(logger)),
        grpc.WithUnaryInterceptor(interceptors.Logging(logger)),
        grpc.WithUnaryInterceptor(interceptors.Metrics()),
    )
    
    pb.RegisterGreeterServer(server, &greeterServer{logger: logger})
    
    server.Start()
}
```

**3. Client usage:**
```go
package main

import (
    "context"
    "fmt"
    "grouter/pkg/messaging/grpc"
    pb "grouter/api/proto/helloworld"
    "time"
)

func main() {
    logger, _ := zap.NewProduction()
    
    config := grpc.ClientConfig{
        Target:  "localhost:9090",
        Timeout: 10 * time.Second,
    }
    
    client, _ := grpc.NewClient(config, logger)
    defer client.Close()
    
    greeter := pb.NewGreeterClient(client.GetConn())
    
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    
    resp, err := greeter.SayHello(ctx, &pb.HelloRequest{Name: "World"})
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    fmt.Println(resp.Message) // Output: Hello World
}
```

---

## Testing

### Unit Testing Interceptors

```go
package myservice_test

import (
    "context"
    "testing"
    "grouter/pkg/messaging/grpc/interceptors"
    "go.uber.org/zap"
    "google.golang.org/grpc"
)

func TestRecoveryInterceptor(t *testing.T) {
    logger := zap.NewNop()
    interceptor := interceptors.Recovery(logger)
    
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        panic("test panic")
    }
    
    info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
    
    _, err := interceptor(context.Background(), nil, info, handler)
    
    if err == nil {
        t.Error("expected error from panic, got nil")
    }
}
```

### Integration Testing

```go
func TestGreeterService(t *testing.T) {
    // Start server
    server := grpc.NewServer(zap.NewNop(), grpc.WithPort(0)) // random port
    pb.RegisterGreeterServer(server, &greeterServer{})
    go server.Start()
    defer server.Stop(context.Background())
    
    // Create client
    client, _ := grpc.NewClient(grpc.ClientConfig{Target: "localhost:9090"}, zap.NewNop())
    defer client.Close()
    
    // Test RPC
    greeter := pb.NewGreeterClient(client.GetConn())
    resp, err := greeter.SayHello(context.Background(), &pb.HelloRequest{Name: "Test"})
    
    assert.NoError(t, err)
    assert.Equal(t, "Hello Test", resp.Message)
}
```

---

## Troubleshooting

### Common Issues

**1. "connection refused"**
```
Error: rpc error: code = Unavailable desc = connection refused
```
Fix: Check server is running and port is correct.

**2. "rate limit exceeded"**
```
Error: rpc error: code = ResourceExhausted desc = rate limit exceeded
```
Fix: Increase rate limiter settings or reduce request rate.

**3. "circuit breaker is open"**
```
Error: rpc error: code = Unavailable desc = circuit breaker is open
```
Fix: Wait for circuit breaker timeout, check downstream service health.

**4. "request timeout exceeded"**
```
Error: rpc error: code = DeadlineExceeded desc = request timeout exceeded
```
Fix: Increase timeout or optimize handler performance.

---

## Summary

The gRPC package provides production-ready server and client implementations with:

- ✅ 7 built-in interceptors
- ✅ Flexible configuration
- ✅ Full observability
- ✅ Client resilience patterns
- ✅ Easy integration

For architecture details, see `ARCHITECTURE.md`.
