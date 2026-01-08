package config

// AppConfig holds service-specific configuration
type AppConfig struct {
	EnableEventPublishing bool `mapstructure:"enable_event_publishing"`
}

// DefaultConfig returns sensible defaults
func DefaultConfig() *AppConfig {
	return &AppConfig{
		EnableEventPublishing: true,
	}
}
