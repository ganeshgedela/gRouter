# Messaging + RPC Service Template

This template combines **NATS messaging** (async pub/sub) with **gRPC** (sync RPC calls) for microservices that need both communication patterns.

## Architecture

```
┌────────────────────────────────────────┐
│   Messaging + RPC Service              │
│                                        │
│  ┌──────────────┐  ┌────────────────┐ │
│  │              │  │                │ │
│  │  NATS        │  │     gRPC       │ │
│  │  Messaging   │  │     Server     │ │
│  │              │  │                │ │
│  │  Async       │  │     Sync       │ │
│  │  Pub/Sub     │  │     RPC        │ │
│  │              │  │                │ │
│  └──────────────┘  └────────────────┘ │
│                                        │
│   Factory Pattern Service Manager      │
└────────────────────────────────────────┘
```

## Features

✅ **NATS Messaging**
- Async pub/sub communication
- Event-driven architecture
- Subscribe to multiple topics
- Message envelope with metadata

✅ **gRPC Server**
- Synchronous RPC calls
- Type-safe proto definitions
- Bidirectional streaming support
- High performance

✅ **Factory Pattern**
- Auto-discovery of message handlers
- Config-driven service loading
- Clean dependency injection

✅ **Use Cases**
- Receive events via NATS, expose status via gRPC
- Process commands async, return results via RPC
- Event sourcing with query API
- CQRS pattern implementation

## Quick Start

### 1. Build the Service
```bash
cd templates/messaging-rpc-service
go build -o bin/server cmd/server/main.go
```

### 2. Run the Service
```bash
./bin/server --config internal/config/config.yaml
```

### 3. Test NATS Messaging
```bash
# Publish event
nats pub events.user.created '{"user_id":"123","name":"John"}'

# Service will log: "received event"
```

### 4. Test gRPC
```bash
# Use grpcurl or client
grpcurl -plaintext -d '{"name":"Alice"}' localhost:9091 hello.HelloService/SayHello
```

## Configuration

```yaml
grpc:
  enabled: true
  port: 9091

nats:
  enabled: true
  url: "nats://localhost:4222"

services:
  message_handler:
    enabled: true
    topics:
      - "events.>"
      - "commands.>"
```

## Default Endpoints

**gRPC:** `localhost:9091`
- `HelloService.SayHello` - Example RPC method

**NATS Topics:**
- `events.>` - All event messages
- `commands.>` - All command messages

## Extending

### Add New gRPC Methods
1. Update `api/proto/hello.proto`
2. Regenerate proto: `protoc --go_out=. --go-grpc_out=. api/proto/hello.proto`
3. Implement in `internal/grpc/server.go`

### Add New NATS Handlers
1. Create handler in `internal/pkg/api/`
2. Register factory in `init()`
3. Subscribe to topics in `Start()`

## When to Use

✅ **Use this template when:**
- You need both async events AND sync queries
- Implementing CQRS pattern
- Event sourcing with API layer
- Processing background tasks with status API

❌ **Don't use when:**
- You only need one protocol (use specific template)
- REST API is preferred over gRPC
- Simple request/response is sufficient

## Related Templates

- `nats-service` - NATS only
- `grpc-service` - gRPC only  
- `hybrid-service` - NATS + REST
- `rest-service` - REST only
