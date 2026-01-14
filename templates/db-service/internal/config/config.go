package config

// AppConfig holds application-specific configuration
type AppConfig struct {
	ServiceName string
	Version     string
}

// DefaultConfig returns default application configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		ServiceName: "db-service",
		Version:     "1.0.0",
	}
}
