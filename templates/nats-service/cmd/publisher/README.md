# NATS Publisher Example

This example demonstrates how to publish messages to test the NATS worker service.

## Usage

### 1. Start the NATS Worker Service
```bash
cd templates/nats-service
go run cmd/worker/main.go --config internal/config/config.yaml
```

### 2. Run the Publisher (in another terminal)
```bash
cd templates/nats-service
go run cmd/publisher/main.go --config internal/config/config.yaml
```

### 3. Observe the Output

**Publisher Output:**
```
🚀 NATS Publisher Example - Testing NATS Worker Service
✅ Connected to NATS server
📤 Test 1: Publishing to user.created topic
✅ Published to user.created
📤 Test 2: Publishing to order.created topic
✅ Published to order.created
📤 Test 3: Publishing to order.updated topic
✅ Published to order.updated
🎉 All test messages published successfully!
```

**Worker Output:**
```
✅ user service subscriptions registered
✅ order service subscriptions registered
🔔 Received user.created event
🔔 Received order.created event
🔔 Received order.updated event
```

## What This Tests

- ✅ NATS connectivity
- ✅ Message publishing
- ✅ Worker subscription handling
- ✅ Service-to-service communication via NATS
- ✅ Factory pattern service discovery

## Topics

The publisher sends messages to these topics that the worker subscribes to:

- `user.created` → User Service
- `order.created` → Order Service
- `order.updated` → Order Service
