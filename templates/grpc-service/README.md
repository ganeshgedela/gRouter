# gRPC Microservice Template

This template demonstrates a **production-grade gRPC service** using the `pkg/messaging/grpc` shared package.

## Architecture

-   **Manager**: Handles lifecycle (Init, Start, Stop) and dependency injection
-   **gRPC Server**: Uses `pkg/messaging/grpc` for high-performance RPC
-   **Protocol Buffers**: Type-safe contract definition
-   **Config**: Centralized configuration management
-   **Logger**: Structured logging with Zap

## Directory Structure

```
.
├── api/
│   └── proto/              # Protocol Buffer definitions
├── cmd/
│   └── server/             # Main entry point
├── internal/
│   ├── app/                # Application lifecycle management
│   │   └── app.go          # Main application orchestrator
│   ├── config/             # Service-specific configuration
│   │   └── config.go
│   └── grpc/               # Service implementation
│       └── server.go
├── Dockerfile
└── README.md
```

## Features

### 1. Production-Ready Patterns
- ✅ Automatic interceptors (Recovery, Logging, Metrics, Tracing)
- ✅ Server reflection for development
- ✅ Graceful shutdown
- ✅ Connection management
- ✅ Health checks

### 2. gRPC Features
- Unary RPC calls
- Streaming support (can be extended)
- Load balancing ready
- TLS/mTLS support (configurable)

## How to Run

### Prerequisites
-   Go 1.22+
-   `protoc` compiler installed
-   `protoc-gen-go` and `protoc-gen-go-grpc` installed

### 1. Generate Protobuf Code
Run from project root:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    templates/grpc-service/api/proto/hello.proto
```

### 2. Run the Server
```bash
go run templates/grpc-service/cmd/server/main.go
```

### 3. Test with grpcurl
```bash
# List services
grpcurl -plaintext localhost:9090 list

# Call SayHello
grpcurl -plaintext -d '{"name": "Developer"}' \
    localhost:9090 hello.HelloService/SayHello
```

Expected response:
```json
{
  "message": "Hello Developer from Production gRPC Service!"
}
```

### Docker Build

```bash
docker build -f templates/grpc-service/Dockerfile -t grpc-server:latest .
docker run -p 9090:9090 grpc-server:latest
```

## Protocol Buffer Definition

The service implements a simple greeting RPC:

```protobuf
service HelloService {
  rpc SayHello (HelloRequest) returns (HelloResponse) {}
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}
```

## Adding New RPCs

1.  Update `api/proto/hello.proto`:
    ```protobuf
    rpc GetUser (GetUserRequest) returns (GetUserResponse) {}
    ```

2.  Regenerate code:
    ```bash
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        templates/grpc-service/api/proto/hello.proto
    ```

3.  Implement in `internal/grpc/server.go`:
    ```go
    func (s *HelloServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
        // Implementation
    }
    ```

## Monitoring

- **Metrics**: Interceptor automatically exports Prometheus metrics
- **Tracing**: OpenTelemetry tracing enabled
- **Logs**: Structured JSON logs with request metadata

## Production Checklist

- ✅ Graceful shutdown with signal handling
- ✅ Connection pooling and keepalive
- ✅ Request/Response logging
- ✅ Panic recovery
- ✅ OpenTelemetry tracing
- ✅ Prometheus metrics
- ✅ Server reflection (disable in production)
