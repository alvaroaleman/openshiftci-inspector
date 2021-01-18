package http

import (
	"crypto/tls"
	"strings"
)

// CreateTLSConfig creates a TLS config. Should only be called after config.Validate().
func (c Config) CreateTLSConfig() (tlsConfig *tls.Config, err error) {
	if !strings.HasPrefix(c.BaseURL, "https://") {
		return nil, nil
	}

	tlsConfig = &tls.Config{
		MinVersion: tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}
	if c.caCertPool != nil {
		tlsConfig.RootCAs = c.caCertPool
	}
	return tlsConfig, nil
}
