package config

import "time"

// ServiceConfig holds service-specific configuration
type ServiceConfig struct {
	Enabled     bool          `mapstructure:"enabled"`
	EnableCache bool          `mapstructure:"enable_cache"`
	CacheTTL    time.Duration `mapstructure:"cache_ttl"`
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
			EnableCache: true,
			CacheTTL:    5 * time.Minute,
		},
		Orders: ServiceConfig{
			Enabled:     true,
			EnableCache: true,
			CacheTTL:    10 * time.Minute,
		},
	}
}
