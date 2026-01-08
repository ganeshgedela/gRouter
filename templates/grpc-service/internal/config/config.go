package config

// AppConfig holds service-specific configuration
type AppConfig struct {
	MaxConcurrentStreams uint32 `mapstructure:"max_concurrent_streams"`
}

// DefaultConfig returns sensible defaults
func DefaultConfig() *AppConfig {
	return &AppConfig{
		MaxConcurrentStreams: 100,
	}
}
