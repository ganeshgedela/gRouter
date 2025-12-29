# Learning Bazel for gRouter

This document provides a basic introduction to Bazel, the build system used in `gRouter`, and how to manage it using Gazelle.

## 1. Basic Tutorial
*   **Workspace**: The root directory of the project, containing a `WORKSPACE` file.
*   **Package**: A directory containing a `BUILD` (or `BUILD.bazel`) file.
*   **Target**: A buildable unit defined in a `BUILD` file.

### Common Commands
*   **Build**: `bazel build //pkg/manager:manager`
*   **Test**: `bazel test //pkg/...`
*   **Clean**: 
    *   `bazel clean` (Removes output directories)
    *   `bazel clean --expunge` (Removes entire working tree and cached state)

### Bazel vs Go Build

| Feature | `go build` | `bazel build` |
| :--- | :--- | :--- |
| **Scope** | Go language only. | Polyglot (Go, C++, Java, Docker, K8s, etc.). |
| **Dependencies** | Implicit via `go.mod`. | Explicit via `BUILD` files and `deps`. |
| **Reproducibility** | Relies on local env (versions, C dependencies). | Hermetic (sandboxed, explicit toolchains). |
| **Caching** | Package-level caching. | Artifact-based, remote caching, aggressive correctness. |
| **Speed** | Fast for small/medium Go projects. | Fastest for large, multi-language monorepos (incremental). |
| **Artifacts** | Binary in current directory. | `bazel-bin/path/to/target`. |

---

## 4. Development & Verification Workflow

### 1. Build
**Using Bazel (Recommended)**
```bash
bin/bazel build //...
```
**Using Go (Standard)**
You can also build the service using standard Go tools, provided you are in the service module directory:
```bash
cd services/natsdemosvc
go build ./cmd/natsdemosvc
# Output binary: natsdemosvc
```

### 2. Unit Tests
**Using Bazel**
```bash
bin/bazel test //...
```
**Using Go**
```bash
cd services/natsdemosvc
go test ./...
```
*Note: Some integration tests may skip if NATS is not running.*

### 3. Integration Tests
1.  **Start NATS Server**:
    ```bash
    docker run -d --name nats-test -p 4222:4222 -p 8222:8222 nats:latest -js
    ```
2.  **Run Tests (Bazel)**:
    ```bash
    bin/bazel test //... --cache_test_results=no
    ```

### 4. Verification Steps
After building, you can manually run and verify the application.

#### A. Run the Application
1.  **Start NATS** (if not running):
    ```bash
    docker run -d --name nats-test -p 4222:4222 -p 8222:8222 nats:latest -js
    ```
2.  **Run the Service**:
    *Via Go*:
    ```bash
    cd services/natsdemosvc
    go run cmd/natsdemosvc/main.go
    ```
    *Via Bazel binary*:
    ```bash
    ./bazel-bin/services/natsdemosvc/cmd/natsdemosvc/natsdemosvc_/natsdemosvc
    ```
    *Running with custom config:*
    ```bash
    ./bazel-bin/services/natsdemosvc/cmd/natsdemosvc/natsdemosvc_/natsdemosvc --config services/natsdemosvc/internal/config/config.yaml
    ```

#### B. Manual Verification (Using NATS Client)
You can use the `nats` CLI tool to interact with the service manually.

1.  **Install NATS CLI**:
    ```bash
    go install github.com/nats-io/natscli/nats@latest
    ```
2.  **Subscribe to Logs/Events**:
    Open a terminal and watch traffic:
    ```bash
    nats sub "natsdemosvc.>"
    ```
3.  **Send Control Commands**:
    In another terminal, send commands (assuming custom config with name `natsdemosvc`):
    *   **Start**:
        ```bash
        nats pub natsdemosvc.start '{"type": "start", "id": "1", "data": {}, "source": "cli"}'
        ```
    *   **Create**:
        ```bash
        nats pub natsdemosvc.natdemo.create '{"type": "natdemo.create", "id": "1", "data": {}, "source": "cli"}'
        ```
    *   **Stop**:
        ```bash
        nats pub natsdemosvc.stop '{"type": "stop", "id": "2", "data": {}, "source": "cli"}'
        ```
4.  **Observe**:
    Check the application logs to confirm "Creating NATS" and other status messages appear.

#### C. Verification Tool
Alternatively, use the included verifier tool:
```bash
cd services/natsdemosvc
go run cmd/verifier/main.go
```
