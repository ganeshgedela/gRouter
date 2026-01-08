package web

import (
	"fmt"
	"os"

	"grouter/pkg/config"
)

// ValidateConfig validates web configuration
func ValidateConfig(cfg config.WebConfig) error {
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", cfg.Port)
	}

	if cfg.ReadTimeout < 0 {
		return fmt.Errorf("read_timeout cannot be negative")
	}

	if cfg.WriteTimeout < 0 {
		return fmt.Errorf("write_timeout cannot be negative")
	}

	if cfg.ShutdownTimeout < 0 {
		return fmt.Errorf("shutdown_timeout cannot be negative")
	}

	// Validate TLS configuration
	if cfg.TLS.Enabled {
		if cfg.TLS.CertFile == "" {
			return fmt.Errorf("TLS enabled but cert_file is empty")
		}
		if cfg.TLS.KeyFile == "" {
			return fmt.Errorf("TLS enabled but key_file is empty")
		}

		// Check if cert and key files exist
		if _, err := os.Stat(cfg.TLS.CertFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS cert file not found: %s", cfg.TLS.CertFile)
		}
		if _, err := os.Stat(cfg.TLS.KeyFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS key file not found: %s", cfg.TLS.KeyFile)
		}
	}

	// Validate mode
	validModes := map[string]bool{
		"debug":      true,
		"release":    true,
		"test":       true,
		"production": true,
	}
	if cfg.Mode != "" && !validModes[cfg.Mode] {
		return fmt.Errorf("invalid mode: %s (must be debug, release, test, or production)", cfg.Mode)
	}

	return nil
}
