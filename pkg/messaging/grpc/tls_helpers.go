package grpc

import "fmt"

// WithTLS is a server option to enable TLS
func WithTLS(certFile, keyFile string) Option {
	return func(s *Server) {
		s.config.TLS = TLSConfig{
			Enabled:  true,
			CertFile: certFile,
			KeyFile:  keyFile,
		}
	}
}

// WithMTLS is a server option to enable mutual TLS (client authentication)
func WithMTLS(certFile, keyFile, caFile string) Option {
	return func(s *Server) {
		s.config.TLS = TLSConfig{
			Enabled:    true,
			CertFile:   certFile,
			KeyFile:    keyFile,
			CAFile:     caFile,
			ClientAuth: true,
		}
	}
}

// WithClientTLS is a client option to enable TLS
func WithClientTLS(serverName string, caFile string) ClientOption {
	return func(c *Client) {
		c.config.TLS = TLSConfig{
			Enabled:    true,
			ServerName: serverName,
			CAFile:     caFile,
		}
	}
}

// WithClientMTLS is a client option to enable mutual TLS
func WithClientMTLS(serverName, certFile, keyFile, caFile string) ClientOption {
	return func(c *Client) {
		c.config.TLS = TLSConfig{
			Enabled:    true,
			ServerName: serverName,
			CertFile:   certFile,
			KeyFile:    keyFile,
			CAFile:     caFile,
		}
	}
}

// WithInsecureTLS is a client option to skip server verification (NOT for production)
func WithInsecureTLS() ClientOption {
	return func(c *Client) {
		c.config.TLS = TLSConfig{
			Enabled:            true,
			InsecureSkipVerify: true,
		}
	}
}

// ValidateTLSConfig validates the TLS configuration
func (c *TLSConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	// Server-side validation
	if c.CertFile != "" || c.KeyFile != "" {
		if c.CertFile == "" {
			return fmt.Errorf("TLS cert file is required when key file is provided")
		}
		if c.KeyFile == "" {
			return fmt.Errorf("TLS key file is required when cert file is provided")
		}
	}

	// Client auth validation
	if c.ClientAuth && c.CAFile == "" {
		return fmt.Errorf("CA file is required when client authentication is enabled")
	}

	return nil
}
