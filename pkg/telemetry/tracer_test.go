package telemetry

import (
	"context"
	"testing"

	"grouter/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestInitTracer(t *testing.T) {
	tests := []struct {
		name      string
		cfg       config.TracingConfig
		expectErr bool
	}{
		{
			name: "Disabled",
			cfg: config.TracingConfig{
				Enabled: false,
			},
			expectErr: false,
		},
		{
			name: "Enabled with stdout exporter",
			cfg: config.TracingConfig{
				Enabled:     true,
				ServiceName: "test-service",
				Exporter:    "stdout",
			},
			expectErr: false,
		},
		{
			name: "Enabled with unknown exporter",
			cfg: config.TracingConfig{
				Enabled:     true,
				ServiceName: "test-service",
				Exporter:    "unknown",
			},
			expectErr: true,
		},
		{
			name: "Enabled with empty exporter (defaults to no-op or stdout logic depending on implementation)",
			cfg: config.TracingConfig{
				Enabled:     true,
				ServiceName: "test-service",
				Exporter:    "",
			},
			expectErr: false, // Currently implementation defaults to no-op if empty string
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shutdown, err := InitTracer(tt.cfg)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, shutdown)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, shutdown)
				// Test shutdown
				err := shutdown(context.Background())
				assert.NoError(t, err)
			}
		})
	}
}
