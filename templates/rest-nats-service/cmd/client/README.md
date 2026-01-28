# REST-NATS Service Client Example

This example demonstrates how to test the hybrid service which combines REST API and NATS messaging.

## Usage

### 1. Start the REST-NATS Server
```bash
cd templates/rest-nats-service
go run cmd/server/main.go --config internal/config/config.yaml
```

### 2. Run the Client (in another terminal)
```bash
cd templates/rest-nats-service
go run cmd/client/main.go --config internal/config/config.yaml
```

### 3. Observe the Output

**Client Output:**
```
🚀 REST-NATS Service Client - Testing REST + NATS
================================================

📋 Test 1: REST Health Checks
  ✅ Liveness probe: OK
  ✅ Readiness probe: {"id":"hybrid-api","name":"REST-NATS API","status":"running"}

🌐 Test 2: REST API Operations
  ⚠️  GET /api/v1/items: Not implemented (expected)
  ⚠️  POST /api/v1/items: Not implemented (expected)

📡 Test 3: NATS Messaging
  ✅ Connected to NATS server
  ✅ Subscribed to test.hybrid topic
  ✅ Published message to test.hybrid
  ✅ Received NATS message: {...}
  ✅ Message roundtrip successful

🎉 All hybrid service tests completed!

The hybrid service successfully supports:
  ✅ REST API endpoints
  ✅ NATS pub/sub messaging
  ✅ Health monitoring
  ✅ Graceful lifecycle management
```

## What This Tests

### REST Capabilities
- ✅ Health check endpoints (`/health/live`, `/health/ready`)
- ✅ RESTful API endpoints (`/api/v1/items`)
- ✅ HTTP request/response handling
- ✅ JSON serialization

### NATS Capabilities
- ✅ NATS connection establishment
- ✅ Topic subscription
- ✅ Message publishing
- ✅ Message receipt and processing
- ✅ Pub/sub roundtrip

## Architecture

The hybrid service combines two communication patterns:

```
┌─────────────────────────────────────┐
│      REST-NATS Service                 │
│                                     │
│  ┌──────────┐      ┌─────────────┐ │
│  │   REST   │      │    NATS     │ │
│  │   API    │      │ Messaging   │ │
│  │          │      │             │ │
│  │ Port     │      │ Pub/Sub     │ │
│  │ 8081     │      │ Events      │ │
│  └──────────┘      └─────────────┘ │
│                                     │
│  Factory Pattern Service Manager    │
└─────────────────────────────────────┘
```

## Use Cases

**REST API:**
- Synchronous request/response
- CRUD operations
- Client-server communication
- HTTP-based integrations

**NATS Messaging:**
- Asynchronous event processing
- Pub/sub pattern
- Microservice communication
- Event-driven architecture

## Building the Client

```bash
go build -o bin/client cmd/client/main.go
./bin/client --config internal/config/config.yaml
```

## Event Publishing

The hybrid service can publish REST events to NATS when `enable_event_publishing: true`:

```yaml
services:
  api:
    enabled: true
    enable_event_publishing: true  # REST events → NATS
```

This allows:
- REST API calls to trigger NATS events
- Other services to react to REST operations
- Event-driven workflows from HTTP requests
