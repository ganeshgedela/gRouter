# Config Package Documentation

The `pkg/config` package provides a centralized, robust, and dynamic configuration management system for the gRouter application. It leverages [Viper](https://github.com/spf13/viper) to handle configuration loading from multiple sources with a defined precedence order.

## Design Philosophy

The configuration system is designed with **loose coupling** and **flexibility** in mind, particularly for the microservices architecture:

1.  **Centralized Core, Decentralized Services**: Core settings (App, NATS, Logging) are strictly typed, while service-specific configurations are handled dynamically (`map[string]interface{}`). This allows services to define their own config structures without modifying the core `config` package.
2.  **Hierarchy & Precedence**: Configuration can be overridden at runtime, facilitating easy switches between development, testing, and production environments.
3.  **Hot Reloading**: The system supports watching for file changes to update configuration without restarting the application.

## Key Features

### 1. Configuration Sources (Precedence Order)
Values are loaded in the following order (highest precedence first):
1.  **Command-Line Flags**: Explicit flags passed at startup (e.g., `--nats-url`).
2.  **Environment Variables**: Variables typically used in Docker/Kubernetes (prefixed with `GROUTER_`).
3.  **Configuration File**: YAML/JSON file (default: `configs/config.yaml`).
4.  **Defaults**: Default values defined in the code.

### 2. Environment Variable Support
*   **Prefix**: `GROUTER_`
*   **Delimiter**: dots (`.`) in keys are replaced by underscores (`_`).
    *   Example: `log.level` becomes `GROUTER_LOG_LEVEL`.

### 3. Dynamic Service Configuration
The `Services` field is defined as `map[string]interface{}`.
*   **Benefit**: This allows the core `gRouter` to load configurations for any arbitrary service module without importing that module's specific types.
*   **Usage**: Services retrieving their config should use `mapstructure` to decode the generic map into their specific struct (e.g., `NATDemoConfig`).

### 4. Configuration Validation
The package includes built-in validation to ensure critical settings are present:
*   **Required Fields**: `app.name`, `nats.url`.
*   **Enum Validation**: `log.level` must be one of `debug`, `info`, `warn`, or `error`.

### 5. Hot Reloading
The `Watch` function utilizes `fsnotify` to monitor the configuration file. When changes are detected:
1.  The file is re-read.
2.  The configuration is re-validated.
3.  A callback function is executed to notify the application of the update.

## Configuration Structure

The core structure is defined in `types.go`:

```go
type Config struct {
    App      AppConfig      // General app info (Name, Version, Env)
    NATS     NATSConfig     // NATS connection details
    Log      LogConfig      // Logging preferences
    Web      WebConfig      // Web server settings
    Services ServicesConfig // Dynamic map for services -> map[string]interface{}
}
```

## Usage Examples

### Initialization and Access
```go
import "grouter/pkg/config"

func main() {
    // 1. Load Configuration
    // Implicitly parses flags like --config configs/production.yaml
    cfg, err := config.Load()
    if err != nil {
        panic(err)
    }

    // 2. Access Global Config Anywhere
    currentCfg := config.Get()
    fmt.Println(currentCfg.App.Name)
}
```

### Watching for Changes
```go
config.Watch(func(newCfg *config.Config) {
    logger.Info("Configuration updated!", zap.String("new_level", newCfg.Log.Level))
    // Re-configure components if necessary
})
```

### Decoding Service Config (Inside a Service)
```go
type MyServiceConfig struct {
    MySetting string `mapstructure:"my_setting"`
}

func NewService() {
    cfg := config.Get()
    var myCfg MyServiceConfig
    
    // Decode the service-specific map into the struct
    if serviceMap, ok := cfg.Services["myservice"]; ok {
        mapstructure.Decode(serviceMap, &myCfg)
    }
}
```

## Command Line Flags
| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to configuration file | `configs/config.yaml` |
| `--log-level` | Override log level | (empty) |
| `--nats-url` | Override NATS URL | (empty) |
