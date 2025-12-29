package config

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	globalConfig *Config
)

// Load initializes and loads configuration from file, environment, and flags
func Load() (*Config, error) {
	// Define command-line flags
	pflag.String("config", "configs/config.yaml", "Path to configuration file")
	pflag.String("log-level", "", "Log level (debug, info, warn, error)")
	pflag.String("nats-url", "", "NATS server URL")
	pflag.Parse()

	// Bind flags to viper
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("failed to bind flags: %w", err)
	}

	// Set config file
	configFile := viper.GetString("config")
	viper.SetConfigFile(configFile)

	// Enable environment variable support
	viper.SetEnvPrefix("GROUTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with command-line flags if provided
	if logLevel := viper.GetString("log-level"); logLevel != "" {
		cfg.Log.Level = logLevel
	}
	if natsURL := viper.GetString("nats-url"); natsURL != "" {
		cfg.NATS.URL = natsURL
	}

	// Validate configuration
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get returns the global configuration
func Get() *Config {
	return globalConfig
}

// Watch watches for configuration changes and reloads
func Watch(callback func(*Config)) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		var cfg Config
		if err := viper.Unmarshal(&cfg); err != nil {
			fmt.Printf("Error reloading config: %v\n", err)
			return
		}
		if err := validate(&cfg); err != nil {
			fmt.Printf("Config validation failed after reload: %v\n", err)
			return
		}
		globalConfig = &cfg
		if callback != nil {
			callback(&cfg)
		}
	})
	viper.WatchConfig()
}

// validate performs configuration validation
func validate(cfg *Config) error {
	if cfg.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}
	if cfg.NATS.Enabled && cfg.NATS.URL == "" {
		return fmt.Errorf("nats.url is required")
	}
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[cfg.Log.Level] {
		return fmt.Errorf("invalid log level: %s", cfg.Log.Level)
	}
	return nil
}
