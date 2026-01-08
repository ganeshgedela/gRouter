package natdemo

// NATDemoConfig holds NATDemo service configuration
type NATDemoConfig struct {
	Enabled bool `mapstructure:"enabled"`
	// Subject prefix for NATS topics
	Subject string `mapstructure:"subject"`
	// User defined name
	Name string `mapstructure:"name"`
	// QueueGroup for load balancing
	QueueGroup string `mapstructure:"queue_group"`
}
