package webdemo

// WebDemoConfig holds WebDemo service configuration
type WebDemoConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Subject string `mapstructure:"subject"` // NATS subject prefix
}
