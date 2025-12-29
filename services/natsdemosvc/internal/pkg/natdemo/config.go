package natdemo

// NATDemoConfig holds NATDemo service configuration
type NATDemoConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Subject string `mapstructure:"subject"` // NATS subject prefix
}
