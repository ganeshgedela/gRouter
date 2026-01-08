package nats

import (
	"fmt"
	"os"
	"testing"
	"time"

	server "github.com/nats-io/nats-server/v2/server"
	natstest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

// RunServer starts an in-process NATS server for testing
func RunServer(opts *server.Options) *server.Server {
	if opts == nil {
		opts = &server.Options{
			Port: -1, // Random port
		}
	}
	return natstest.RunServer(opts)
}

// getNATSURLHelper returns the NATS server URL for tests
func getNATSURL(s *server.Server) string {
	if s != nil {
		return s.ClientURL()
	}
	// Fallback environment variable
	url := os.Getenv("NATS_URL")
	if url == "" {
		fmt.Println("DEBUG: NATS_URL env var is empty, defaulting to localhost")
		return "nats://localhost:4222"
	}
	// fmt.Printf("DEBUG: getNATSURL returning %s\n", url)
	return url
}

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Start a global embedded NATS server for integration tests
	opts := &server.Options{
		Port:      -1, // Random port
		JetStream: true,
	}
	s := natstest.RunServer(opts)
	if s == nil {
		panic("Failed to start NATS server in TestMain")
	}
	defer s.Shutdown()

	url := s.ClientURL()
	fmt.Printf("DEBUG: Global TestMain Server started at %s\n", url)

	// Set NATS_URL environment variable so tests can pick it up if they use os.Getenv
	os.Setenv("NATS_URL", url)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// waitForNATS waits for NATS server to be available (Deprecated with embedded server)
func waitForNATS() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	opts := []nats.Option{nats.Timeout(1 * time.Second)}
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		nc, err := nats.Connect(natsURL, opts...)
		if err == nil {
			nc.Close()
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// skipIfNoNATS skips the test if NATS server is not available (Deprecated with embedded server)
func skipIfNoNATS(t *testing.T) {
	// No-op if we are using embedded server logic in tests
}
