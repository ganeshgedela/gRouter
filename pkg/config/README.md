# Config Package (`pkg/config`)

The `config` package provides a centralized, type-safe configuration management system for the gRouter application. It supports loading configuration from files, environment variables, and command-line flags, with built-in validation and hot-reloading capabilities.

## Features

- **Hierarchical Loading**: Loads config with precedence: Flags > Env Vars > Config File > Defaults.
- **Environment Support**: Automatically binds `GROUTER_` prefixed environment variables (e.g., `GROUTER_NATS_URL`).
- **Hot Reloading**: Watches the configuration file for changes and updates runtime config dynamically.
- **Type Safety**: Unmarshals configuration into structured Go types.

## Usage

### 1. Loading Configuration

```go
import "github.com/ganesh/grouter/pkg/config"

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    fmt.Printf("App Name: %s\n", cfg.App.Name)
}
```

### 2. Watching for Changes

```go
config.Watch(func(newCfg *config.Config) {
    log.Println("Configuration reloaded!")
    // Apply changes (e.g., update log level)
})
```

## Configuration Structure

The configuration is defined in `configs/config.yaml` (default).

```yaml
app:
  name: "gRouter"
  version: "1.0.0"
  environment: "dev"

nats:
  url: "nats://localhost:4222"
  max_reconnects: 5
  reconnect_wait: 2s

log:
  level: "info"
  format: "json" # json or console
  output_path: "stdout"

services:
  ipsec:
    enabled: true
    subject: "grouter.ipsec"
```

## Environment Variables

All keys can be overridden using environment variables with the `GROUTER_` prefix. Dots are replaced by underscores.

*   `app.name` -> `GROUTER_APP_NAME`
*   `nats.url` -> `GROUTER_NATS_URL`
*   `log.level` -> `GROUTER_LOG_LEVEL`
