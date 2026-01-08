package nats

import (
	"os"

	server "github.com/nats-io/nats-server/v2/server"
	natstest "github.com/nats-io/nats-server/v2/test"
)

// RunTestServer starts an in-process NATS server for testing
// Renamed from RunTestServer to avoid conflicts
func RunTestServer(opts *server.Options) *server.Server {
	if opts == nil {
		opts = &server.Options{
			Port: -1, // Random port
		}
	}
	return natstest.RunServer(opts)
}

// GetTestNATSURL returns the NATS server URL for tests
// Renamed from GetTestNATSURL to avoid conflicts and export for tests
func GetTestNATSURL(s *server.Server) string {
	if s != nil {
		return s.ClientURL()
	}
	// Fallback environment variable
	url := os.Getenv("NATS_URL")
	if url == "" {
		// fmt.Println("DEBUG: NATS_URL env var is empty, defaulting to localhost")
		return "nats://localhost:4222"
	}
	return url
}

// TestMainWrapper can be called by individual test files if needed,
// but for now we rely on the implementation in this file to provide the environment.
// Since TestMain was reported as redefined, we won't define it here effectively ensuring
// only one TestMain exists (whichever one the user has hidden).
// However, to ensure NATS_URL is set if no other TestMain exists, we provide a helper to init env.

func init() {
	// If no NATS_URL is present, we set default for local development/testing convenience
	if os.Getenv("NATS_URL") == "" {
		os.Setenv("NATS_URL", "nats://localhost:4222")
	}
}
