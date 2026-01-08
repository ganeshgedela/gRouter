package config

import "time"

// ServiceConfig holds service-specific configuration
type ServiceConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	ServiceName string        `mapstructure:"service_name"`
	WorkerCount int           `mapstructure:"worker_count"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

// AppConfig extends the global config with service-specific settings
type AppConfig struct {
	Users  ServiceConfig `mapstructure:"users"`
	Orders ServiceConfig `mapstructure:"orders"`
}

// DefaultConfig returns sensible defaults
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Users: ServiceConfig{
			Enabled:     true,
			ServiceName: "user-service",
			WorkerCount: 3,
			Timeout:     10 * time.Second,
		},
		Orders: ServiceConfig{
			Enabled:     true,
			ServiceName: "order-service",
			WorkerCount: 5,
			Timeout:     30 * time.Second,
		},
	}
}
