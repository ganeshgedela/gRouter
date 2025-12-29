# Standard Web Microservice Architecture

For a production-ready web microservice in the `gRouter` ecosystem, the following "internal services" (components) are mandatory to ensure observability, reliability, and maintainability.

## 1. Mandatory Services (Internal Components)

1.  **Configuration Service** (`pkg/config`)
    *   **Role**: Loads environment variables, config files, and flags.
    *   **Why**: 12-factor app compliance; separate code from config.
2.  **Structured Logging** (`pkg/logger`)
    *   **Role**: JSON-structured logs with correlation IDs.
    *   **Why**: Debugging in production (Splunk/ELK).
3.  **Health Service** (`pkg/health`)
    *   **Role**: Liveness (process up) and Readiness (can take traffic) checks.
    *   **Why**: Required by Orchestrators (Kubernetes) for zero-downtime deployments.
4.  **Web Server / Router** (`pkg/web`)
    *   **Role**: HTTP handling, middleware (Auth, CORS, RequestID).
    *   **Why**: Core business logic exposure.
5.  **Metrics / Instrumentation** (Prometheus in `pkg/web`)
    *   **Role**: RED method metrics (Rate, Errors, Duration).
    *   **Why**: Monitoring operational health and performance.
6.  **Signal Handling / Graceful Shutdown**
    *   **Role**: Catch SIGTERM/SIGINT.
    *   **Why**: Finish in-flight requests before exiting; prevent data loss.

## 2. Implementation Sketch

This sketch demonstrates how to wire these mandatory services together using a dedicated `App` struct.

### Directory Structure
```text
internal/
  app/
    app.go       // Wiring logic
    handlers.go  // Business logic
  config/
    config.yaml  // Default config
cmd/
  myservice/
    main.go      // Entry point
```

### Code Sketch (`internal/app/app.go`)

```go
package app

import (
    "context"
    "fmt"
    "go.uber.org/zap"
    
    "grouter/pkg/config"
    "grouter/pkg/health"
    "grouter/pkg/logger"
    "grouter/pkg/web"
)

// App holds all mandatory services
type App struct {
    Cfg    *config.Config
    Log    *zap.Logger
    Health *health.HealthService
    Server *web.Server
}

// New creates the app container (Dependency Injection root)
func New() *App {
    return &App{
        Health: health.NewHealthService(),
    }
}

// Init initializes all mandatory services in the correct order
func (a *App) Init() error {
    // 1. Config
    cfg, err := config.Load()
    if err != nil { return err }
    a.Cfg = cfg

    // 2. Logging
    log, err := logger.New(logger.Config{Level: cfg.Log.Level})
    if err != nil { return err }
    a.Log = log
    
    // 3. Health Checks
    a.Health.AddLivenessCheck("runtime", func() error { return nil })
    // Add other checks (DB, Redis) here

    // 4. Web Server
    a.Server = web.NewWebServer(web.Config{
        Port: cfg.Web.Port,
        Metrics: web.MetricsConfig{Enabled: true}, // 5. Metrics
    }, a.Log, a.Health)

    // Register Business Routes
    a.Server.RegisterService(a)

    return nil
}

// RegisterRoutes implements web.Service
func (a *App) RegisterRoutes(r *gin.RouterGroup) {
    r.GET("/api/v1/resource", a.HandleResource)
}

// Start runs the server
func (a *App) Start(ctx context.Context) error {
    a.Log.Info("Starting service...")
    return a.Server.Start()
}

// Stop handles graceful shutdown
func (a *App) Stop(ctx context.Context) error {
    a.Log.Info("Stopping service...")
    return a.Server.Stop(ctx)
}
```

### Entry Point (`main.go`)

```go
func main() {
    app := app.New()
    if err := app.Init(); err != nil { panic(err) }

    // block until signal
    // ... signal handling code ...
    // app.Stop(ctx)
}
```
