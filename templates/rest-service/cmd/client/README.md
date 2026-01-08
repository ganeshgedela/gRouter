# REST API Client Example

This example demonstrates how to test the REST API service endpoints.

## Usage

### 1. Start the REST API Service
```bash
cd templates/rest-service
go run cmd/api/main.go --config internal/config/config.yaml
```

### 2. Run the Client (in another terminal)
```bash
cd templates/rest-service
go run cmd/client/main.go
```

### 3. Observe the Output

**Client Output:**
```
🚀 REST API Client Example - Testing REST API Service
==============================================================

📤 Test 1: Health Check
✅ GET /api/v1/health
   Status: 200 OK
   Response: {
     "id": "rest-api",
     "name": "REST API Service",
     "status": "alive"
   }

📤 Test 2: Liveness Probe
✅ GET /health/live
   Status: 200 OK
   ...

📤 Test 3: Readiness Probe
✅ GET /health/ready
   Status: 200 OK
   ...

🎉 All API tests completed!
```

## What This Tests

- ✅ Health check endpoints
- ✅ Liveness probe
- ✅ Readiness probe
- ✅ User API endpoints (example)
- ✅ Order API endpoints (example)
- ✅ JSON request/response handling
- ✅ HTTP status codes

## Available Endpoints

The client tests these endpoints:

**Health & Monitoring:**
- `GET /api/v1/health` → Application health status
- `GET /health/live` → Liveness probe
- `GET /health/ready` → Readiness probe

**User Service:**
- `POST /api/v1/users` → Create user
- `GET /api/v1/users` → List users

**Order Service:**
- `POST /api/v1/orders` → Create order
- `GET /api/v1/orders` → List orders

## Building the Client

```bash
go build -o bin/client cmd/client/main.go
./bin/client
```
