# WebDemoSvc

A pure HTTP service demonstrating the use of `pkg/web` within the gRouter ecosystem.

## Prerequisites

The service currently depends on NATS for the `ServiceManager`. You must have a NATS server running.

**Start NATS using Docker:**
```bash
docker run --rm --name nats-test -p 4222:4222 -d nats:latest
```

## Option 1: Run Locally (Go)

1.  **Ensure NATS is running** (see Prerequisites).
2.  **Navigate to the service directory:**
    ```bash
    cd services/webdemosvc
    ```
3.  **Run the service:**
    ```bash
    go run cmd/webdemosvc/main.go
    ```
    *Note: By default, it looks for `config.yaml` in the current directory or `internal/config/`. The code uses `services/webdemosvc/internal/config/config.yaml` relative to module root if running from root.*

    **Better way from project root:**
    ```bash
    go run services/webdemosvc/cmd/webdemosvc/main.go
    ```

## Option 2: Run with Docker

1.  **Ensure NATS is running** (see Prerequisites).
2.  **Build the image:**
    ```bash
    docker build -t webdemosvc:latest -f services/webdemosvc/Dockerfile .
    ```
3.  **Run the container:**
    *We use `--network host` to easily connect to the local NATS server running on port 4222.*
    ```bash
    docker run --rm --network host --name webdemosvc-run webdemosvc:latest
    ```

## Verification

Once running (default port `8080`), you can test the endpoints:

**1. Hello Endpoint**
```bash
curl http://localhost:8080/hello
# Output: {"message":"Hello from WebDemoSvc!"}
```

**2. Echo Endpoint**
```bash
curl "http://localhost:8080/echo?msg=gRouter"
# Output: {"echo":"gRouter"}
```

**3. Health Check**
```bash
curl http://localhost:8080/health/live
# Output: {"checks":{},"status":"up"}
```
