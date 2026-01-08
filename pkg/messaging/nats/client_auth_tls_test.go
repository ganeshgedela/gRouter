package nats

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Tests for Configuration logic using the internal buildConnectionOptions method.
// This gives us high coverage of the configuration logic without needing a running NATS server.

func TestClient_TLSConfig(t *testing.T) {
	logger := zap.NewNop()

	t.Run("TLS_Enabled_SkipVerify", func(t *testing.T) {
		cfg := Config{
			URL: "nats://localhost:4222",
			TLS: TLSConfig{
				Enabled:    true,
				SkipVerify: true,
			},
		}

		client, err := NewNATSClient(cfg, logger)
		assert.NoError(t, err)

		// Verify options are built (len > base options)
		// Base options are ~6 (removed ClosedHandler). TLS adds at least 1 (Secure).
		opts := client.buildConnectionOptions()
		assert.NotEmpty(t, opts)
		// We expect at least base (6) + secure (1) = 7
		assert.GreaterOrEqual(t, len(opts), 7)
	})

	t.Run("TLS_With_Certs", func(t *testing.T) {
		certFile, _ := os.CreateTemp("", "cert")
		keyFile, _ := os.CreateTemp("", "key")
		caFile, _ := os.CreateTemp("", "ca")
		defer os.Remove(certFile.Name())
		defer os.Remove(keyFile.Name())
		defer os.Remove(caFile.Name())

		cfg := Config{
			URL: "nats://localhost:4222",
			TLS: TLSConfig{
				Enabled:  true,
				CertFile: certFile.Name(),
				KeyFile:  keyFile.Name(),
				CAFile:   caFile.Name(),
			},
		}

		client, err := NewNATSClient(cfg, logger)
		assert.NoError(t, err)

		opts := client.buildConnectionOptions()
		// Base (6) + RootCA+ClientCert+Secure (3) = 9
		assert.GreaterOrEqual(t, len(opts), 9)
	})
}

func TestClient_AuthConfig(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Auth_Token", func(t *testing.T) {
		cfg := Config{
			URL: "nats://localhost:4222",
			Auth: AuthConfig{
				Token: "secret_token",
			},
		}
		client, err := NewNATSClient(cfg, logger)
		assert.NoError(t, err)

		opts := client.buildConnectionOptions()
		// Base (6) + Token (1) = 7
		assert.GreaterOrEqual(t, len(opts), 7)
	})

	t.Run("Auth_UserPass", func(t *testing.T) {
		cfg := Config{
			URL: "nats://localhost:4222",
			Auth: AuthConfig{
				Username: "user",
				Password: "password",
			},
		}
		client, err := NewNATSClient(cfg, logger)
		assert.NoError(t, err)

		opts := client.buildConnectionOptions()
		// Base (6) + UserInfo (1) = 7
		assert.GreaterOrEqual(t, len(opts), 7)
	})

	t.Run("Auth_CredsFile", func(t *testing.T) {
		credsFile, _ := os.CreateTemp("", "creds")
		defer os.Remove(credsFile.Name())

		cfg := Config{
			URL: "nats://localhost:4222",
			Auth: AuthConfig{
				CredsFile: credsFile.Name(),
			},
		}
		client, err := NewNATSClient(cfg, logger)
		assert.NoError(t, err)

		opts := client.buildConnectionOptions()
		// Base (6) + Creds (1) = 7
		assert.GreaterOrEqual(t, len(opts), 7)
	})
}
