# Logger Package (`pkg/logger`)

The `logger` package provides a high-performance, structured logging facility for gRouter, built on top of [uber-go/zap](https://github.com/uber-go/zap). It supports context-aware logging, enabling distributed tracing and request-scoped fields.

## Features

- **Structured Logging**: JSON output for production, colored console output for development.
- **Context Awareness**: Propagate logger instances with context-specific fields (e.g., `request_id`, `trace_id`).
- **Performance**: Zero-allocation logging path where possible.
- **Global & Local**: Access via a thread-safe global singleton or local instances.

## Usage

### 1. Initialization

```go
import "github.com/ganesh/grouter/pkg/logger"

conf := logger.Config{
    Level:      "info",
    Format:     "json",
    OutputPath: "stdout",
}

log, err := logger.New(conf)
if err != nil {
    panic(err)
}
// log is now the global logger
```

### 2. Basic Logging

```go
logger.Info("Application started", zap.String("version", "1.0.0"))
logger.Error("Database connection failed", zap.Error(err))
```

### 3. Context-Aware Logging

To trace requests across boundaries, embed the logger in the context:

```go
// Add request ID to context
ctx = logger.WithRequestID(ctx, "req-12345")

// Retrieve and log (automatically includes request_id)
logger.FromContext(ctx).Info("Processing payment")
// Output: {"level":"info","msg":"Processing payment","request_id":"req-12345",...}
```

## Configuration

| Field | Type | Description |
| :--- | :--- | :--- |
| `Level` | `string` | Log level (`debug`, `info`, `warn`, `error`, `fatal`). |
| `Format` | `string` | Output format: `json` (default) or `console`. |
| `OutputPath` | `string` | File path or `stdout`/`stderr`. |
