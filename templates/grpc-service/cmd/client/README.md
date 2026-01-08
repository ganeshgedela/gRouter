# gRPC Client Example

This example demonstrates how to test the gRPC service using a Go client.

## Usage

### 1. Start the gRPC Server
```bash
cd templates/grpc-service
go run cmd/server/main.go --config internal/config/config.yaml
```

### 2. Run the Client (in another terminal)
```bash
cd templates/grpc-service
go run cmd/client/main.go
```

### 3. Observe the Output

**Client Output:**
```
🚀 gRPC Client Example - Testing gRPC Service
============================================================

📡 Connecting to gRPC server at localhost:9090...
✅ Connected to gRPC server

📤 Test 1: Sending SayHello request...
✅ Request: {name: "Alice"}
✅ Response: {message: "Hello, Alice from gRouter Production gRPC Service!"}

📤 Test 2: Sending SayHello request...
✅ Request: {name: "Bob"}
✅ Response: {message: "Hello, Bob from gRouter Production gRPC Service!"}

📤 Test 3: Testing with context timeout...
✅ Request: {name: "Charlie"} (timeout: 2s)
✅ Response: {message: "Hello, Charlie from gRouter Production gRPC Service!"}
⏱️  Duration: 125µs

🎉 All gRPC tests completed successfully!
```

**Server Output:**
```
{"level":"info","msg":"received hello request","name":"Alice"}
{"level":"info","msg":"received hello request","name":"Bob"}
{"level":"info","msg":"received hello request","name":"Charlie"}
```

## What This Tests

- ✅ gRPC connection establishment
- ✅ HelloService.SayHello RPC method
- ✅ Request/response serialization
- ✅ Context and timeout handling
- ✅ Server-side logging

## Proto Definition

The client uses the HelloService defined in `api/proto/hello.proto`:

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

## Building the Client

```bash
go build -o bin/client cmd/client/main.go
./bin/client
```
