package web

import (
	"crypto/tls"
	"fmt"

	"grouter/pkg/config"
)

// LoadTLSConfig creates a TLS configuration from config
func LoadTLSConfig(cfg config.TLSConfig) (*tls.Config, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: secureCipherSuites(),
	}

	// TODO: Add client auth support when config fields are added
	// For now, TLS is server-side only

	return tlsConfig, nil
}

// setupClientAuth configures client certificate authentication
func setupClientAuth(tlsConfig *tls.Config, cfg config.TLSConfig) error {
	// For now, we don't have client auth fields in config
	// This will be a future enhancement
	return nil
}

// secureCipherSuites returns a list of secure cipher suites
func secureCipherSuites() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	}
}
