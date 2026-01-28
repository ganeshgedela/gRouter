package config

// AppConfig represents application-specific configuration
type AppConfig struct {
	// Add any app-specific config here
}

// DefaultConfig returns the default application configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{}
}
