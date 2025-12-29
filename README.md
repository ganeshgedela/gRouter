# gRouter - Event-Driven Microservice

A Go-based event-driven microservice architecture using NATS for message handling. Supports multiple independent services (IPSec, interface, firewall, etc.) with a modular package structure.

## Features

- **Event-Driven Architecture**: Built on NATS for scalable, asynchronous messaging
- **Modular Services**: Independent service packages that can be enabled/disabled
- **Configuration Management**: Viper + pflag for flexible config from files, env vars, and CLI flags
- **Structured Logging**: Zap logger with JSON and console output formats
- **Bazel Build System**: Hermetic, incremental builds
- **Docker Support**: Multi-stage Dockerfile for minimal production images

## Project Structure

### Source Code Application
*   **`services/`**: Contains the source code for independent microservices (e.g., `natsdemosvc`).
*   **`pkg/`**: Shared library code used across services.
    *   `config/`: Configuration management (Viper + pflag).
    *   `logger/`: Structured logging framework (Zap).
    *   `messaging/`: NATS event handling and client wrappers.
    *   `manager/`: Service orchestration logic.
*   **`cmd/`**: Entry points for older or monolithic applications (legacy).
*   **`configs/`**: Runtime configuration files (e.g., `config.yaml`).
*   **`api/`**: API definitions (Protobufs, OpenAPI/Swagger specs).
*   **`deployments/`**: Deployment assets (Dockerfiles, Kubernetes manifests).
*   **`docs/`**: Generated documentation (Swagger UI files).

### Build System (Bazel)
*   **`WORKSPACE`**: Root definition of external dependencies.
*   **`BUILD.bazel`**: Root build configuration and Gazelle settings.
*   **`MODULE.bazel` / `MODULE.bazel.lock`**: Bzlmod dependency management.
*   **`deps.bzl`**: Generated list of Go repositories for Bazel.
*   **`bazel-*`**: Symlinks to build artifacts and output directories (e.g., `bazel-bin` for compiled binaries).
*   **`bin/`**: Helper scripts or local binaries (e.g., local `bazel` wrapper).

### Go Language
*   **`go.mod` / `go.sum`**: Go module definitions and checksums.

### Documentation
*   **`production_readiness.md`**: Checklist for production requirements.
*   **`learning_bazel.md`**: Guide for building/testing with Bazel.
*   **`app.log`**: Runtime logs (if enabled).

## Prerequisites

### For Go Build
- Go 1.22 or higher
- NATS server (for runtime)

### For Bazel Build
- Bazel 6.0+ or Bazelisk
- Go 1.22 or higher

## Quick Start

### 1. Install Dependencies

```bash
go mod download
```

### 2. Start NATS Server

```bash
# Using Docker
docker run -d --name nats -p 4222:4222 nats:latest

# Or install locally
# https://docs.nats.io/running-a-nats-service/introduction/installation
```

### 3. Build and Run

#### Using Go

```bash
# Build
go build -o bin/service ./cmd/service

# Run
./bin/service --config configs/config.yaml
```

#### Using Bazel

```bash
# Install Bazel (if not already installed)
# Option 1: Install Bazelisk (recommended)
curl -Lo /usr/local/bin/bazel https://github.com/bazelbuild/bazelisk/releases/download/v1.19.0/bazelisk-linux-amd64
chmod +x /usr/local/bin/bazel

# Option 2: Install via package manager
sudo apt install bazel-bootstrap  # Ubuntu/Debian

# Build
bazel build //cmd/service:service

# Run
bazel run //cmd/service:service -- --config configs/config.yaml
```

#### Using Docker

**Important**: You must run the build command from the **root of the repository** to ensure all dependencies (like `go.mod`) are correctly included.

```bash
# 1. Build the image (Example: natsdemosvc)
docker build -f services/natsdemosvc/Dockerfile -t natsdemosvc:latest .

# 2. Start NATS Server (if not running)
docker run -d --name nats-test -p 4222:4222 -p 8222:8222 nats:latest -js

# 3. Run the Service
# Note: --network host is used here for simplicity to connect to localhost NATS.
# In production, use a Docker network.
docker run -d \
  --name natsdemosvc \
  --network host \
  natsdemosvc:latest
```

### 4. Verification (Docker)

You can verify the service using the NATS CLI. If you don't have it installed locally, you can use the `nats-box` Docker image.

**A. Subscribe to Logs**
```bash
docker logs -f natsdemosvc
```

**B. Send Start Signal**
```bash
docker run --rm --network host natsio/nats-box \
  nats pub natsdemosvc.start '{"type": "start", "id": "1", "data": {}, "source": "cli"}'
```
*Expected Log Output:* `Start signal received. Registering services...`

**C. Send Create Signal**
```bash
docker run --rm --network host natsio/nats-box \
  nats pub natsdemosvc.natdemo.create '{"type": "natdemo.create", "id": "2", "data": {}, "source": "cli"}'
```
*Expected Log Output:* `Creating NATS`

## Configuration

Configuration can be provided via:
1. YAML file (default: `configs/config.yaml`)
2. Environment variables (prefix: `GROUTER_`)
3. Command-line flags

### Example Configuration

```yaml
app:
  name: "gRouter"
  version: "1.0.0"
  environment: "development"

nats:
  url: "nats://localhost:4222"
  max_reconnects: 10
  reconnect_wait: 2s
  connection_timeout: 5s

log:
  level: "info"
  format: "console"
  output_path: "stdout"

services:
  ipsec:
    enabled: true
    subject: "ipsec"
```

### Command-Line Flags

```bash
./bin/service \
  --config configs/config.yaml \
  --log-level debug \
  --nats-url nats://localhost:4222
```

### Environment Variables

```bash
export GROUTER_LOG_LEVEL=debug
export GROUTER_NATS_URL=nats://localhost:4222
./bin/service
```

## Usage

### IPSec Service

The IPSec service listens for events on NATS subjects and manages IPSec tunnels.

#### Create Tunnel

```bash
# Using NATS CLI
nats pub ipsec.tunnel.create '{
  "name": "tunnel-1",
  "source": "10.0.0.1",
  "destination": "10.0.0.2",
  "local_subnet": "192.168.1.0/24",
  "remote_subnet": "192.168.2.0/24",
  "pre_shared_key": "secret",
  "encryption_algo": "aes256",
  "auth_algo": "sha256",
  "dpd_interval": 30
}'
```

#### Delete Tunnel

```bash
nats pub ipsec.tunnel.delete '{"id": "tunnel-id"}'
```

#### Get Status

```bash
nats pub ipsec.tunnel.status '{"id": "tunnel-id"}'
```

## Development

### Running Tests

```bash
# Using Go
go test ./pkg/... -v

# Using Bazel
bazel test //pkg/...

# With coverage
bazel coverage //pkg/...
```

### Adding a New Service

1. Create service package: `pkg/services/myservice/`
2. Implement models, events, service logic, and handlers
3. Add configuration in `pkg/config/types.go`
4. Wire up in `cmd/service/main.go`
5. Update `configs/config.yaml`
6. Create `BUILD.bazel` file

## Architecture

### Event Flow

```
Client → NATS → Subscriber → Handler → Service → Response → NATS → Client
```

### Components

- **Config**: Centralized configuration management
- **Logger**: Structured logging with context propagation
- **Messaging**: NATS client, publisher, and subscriber abstractions
- **Services**: Independent service implementations (IPSec, etc.)

## Production Deployment

### Docker Compose

```yaml
version: '3.8'
services:
  nats:
    image: nats:latest
    ports:
      - "4222:4222"
  
  grouter:
    image: grouter-service:latest
    depends_on:
      - nats
    environment:
      - GROUTER_NATS_URL=nats://nats:4222
      - GROUTER_LOG_LEVEL=info
    volumes:
      - ./configs:/app/configs
```

### Kubernetes

See `deployments/k8s/` (to be added) for Kubernetes manifests.

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
