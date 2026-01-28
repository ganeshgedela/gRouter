package config

// Config holds service configuration
type Config struct {
	HTTPPort string
	GRPCPort string
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		HTTPPort: "8080",
		GRPCPort: "50051",
	}
}
