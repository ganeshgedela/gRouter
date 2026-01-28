# REST-gRPC Service Template

This template demonstrates a service that exposes **both REST (HTTP)** and **gRPC** interfaces, sharing the same business logic.

## Features

- **Dual Protocol Support**:
  - HTTP REST API on port `8080` (Gin framework)
  - gRPC API on port `50051` (grpc-go)
- **Shared Business Logic**: Both interfaces call the same internal service layer
- **Clean Architecture**: Separation of concerns (API -> Business -> Data)
- **Production Ready**: Graceful shutdown, structured logging, health checks

## Directory Structure

```
templates/rest-grpc-service/
├── api/
│   └── proto/           # Protobuf definitions
├── cmd/
│   └── server/          # Main entry point
├── internal/
│   ├── app/             # Application orchestrator
│   ├── config/          # Configuration
│   ├── grpc/            # gRPC server implementation
│   └── pkg/
│       ├── business/    # Business logic (shared)
│       └── rest/        # REST API implementation
```

## Running the Service

```bash
# Run with Bazel
bazel run //templates/rest-grpc-service/cmd/server

# Run with Go
go run templates/rest-grpc-service/cmd/server/main.go
```

## Testing Endpoints

### REST API

```bash
# Health Check
curl http://localhost:8080/api/v1/health

# Create Item
curl -X POST http://localhost:8080/api/v1/items \
  -d '{"name":"Laptop", "description":"Pro model", "price":1999.99, "quantity":10}'

# List Items
curl http://localhost:8080/api/v1/items
```

### gRPC API

Using `grpcurl`:

```bash
# List methods
grpcurl -plaintext localhost:50051 list

# Create Item
grpcurl -plaintext -d '{"name": "Mouse", "price": 29.99, "quantity": 100}' \
  localhost:50051 items.ItemService/CreateItem

# List Items
grpcurl -plaintext localhost:50051 items.ItemService/ListItems
```
