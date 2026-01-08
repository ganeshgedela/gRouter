package web

import (
	"testing"

	"grouter/pkg/config"
)

func TestLoadTLSConfig_Disabled(t *testing.T) {
	cfg := config.TLSConfig{
		Enabled: false,
	}

	tlsConfig, err := LoadTLSConfig(cfg)
	if err != nil {
		t.Errorf("LoadTLSConfig() error = %v, want nil", err)
	}
	if tlsConfig != nil {
		t.Error("Expected nil TLS config when disabled")
	}
}

func TestLoadTLSConfig_MissingFiles(t *testing.T) {
	cfg := config.TLSConfig{
		Enabled:  true,
		CertFile: "/nonexistent/cert.pem",
		KeyFile:  "/nonexistent/key.pem",
	}

	_, err := LoadTLSConfig(cfg)
	if err == nil {
		t.Error("Expected error for missing certificate files")
	}
}

func TestSecureCipherSuites(t *testing.T) {
	suites := secureCipherSuites()

	if len(suites) == 0 {
		t.Error("Expected non-empty cipher suites")
	}

	// Verify all suites are valid
	for _, suite := range suites {
		if suite == 0 {
			t.Error("Invalid cipher suite ID: 0")
		}
	}
}
