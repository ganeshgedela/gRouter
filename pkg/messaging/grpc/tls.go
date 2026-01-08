package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

// configureTLS creates TLS credentials for the server
func configureTLS(config TLSConfig) (credentials.TransportCredentials, error) {
	if !config.Enabled {
		return nil, nil
	}

	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server key pair: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12, // Minimum TLS 1.2 for security
	}

	// Configure client authentication (mTLS)
	if config.ClientAuth && config.CAFile != "" {
		// Load CA certificate
		caCert, err := os.ReadFile(config.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certificate")
		}

		tlsConfig.ClientCAs = certPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return credentials.NewTLS(tlsConfig), nil
}

// configureClientTLS creates TLS credentials for the client
func configureClientTLS(config TLSConfig) (credentials.TransportCredentials, error) {
	if !config.Enabled {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Configure server name for verification
	if config.ServerName != "" {
		tlsConfig.ServerName = config.ServerName
	}

	// Skip verification (NOT recommended for production)
	if config.InsecureSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	// Load CA certificate for server verification
	if config.CAFile != "" {
		caCert, err := os.ReadFile(config.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certificate")
		}

		tlsConfig.RootCAs = certPool
	}

	// Load client certificate for mTLS
	if config.CertFile != "" && config.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client key pair: %w", err)
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return credentials.NewTLS(tlsConfig), nil
}
