package web

import (
	"testing"
	"time"

	"grouter/pkg/config"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.WebConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: config.WebConfig{
				Port:            8080,
				Mode:            "release",
				ReadTimeout:     30 * time.Second,
				WriteTimeout:    30 * time.Second,
				ShutdownTimeout: 10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid port - too low",
			cfg: config.WebConfig{
				Port: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			cfg: config.WebConfig{
				Port: 70000,
			},
			wantErr: true,
		},
		{
			name: "negative read timeout",
			cfg: config.WebConfig{
				Port:        8080,
				ReadTimeout: -5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "negative write timeout",
			cfg: config.WebConfig{
				Port:         8080,
				WriteTimeout: -5 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "TLS enabled without cert",
			cfg: config.WebConfig{
				Port: 8080,
				TLS: config.TLSConfig{
					Enabled: true,
					KeyFile: "/path/to/key.pem",
				},
			},
			wantErr: true,
		},
		{
			name: "TLS enabled without key",
			cfg: config.WebConfig{
				Port: 8080,
				TLS: config.TLSConfig{
					Enabled:  true,
					CertFile: "/path/to/cert.pem",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid mode",
			cfg: config.WebConfig{
				Port: 8080,
				Mode: "invalid-mode",
			},
			wantErr: true,
		},
		{
			name: "valid modes",
			cfg: config.WebConfig{
				Port: 8080,
				Mode: "production",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
