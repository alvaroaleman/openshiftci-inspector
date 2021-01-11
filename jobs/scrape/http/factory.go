package http

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/janoszen/openshiftci-inspector/jobs/scrape"
)

// NewHTTPScraper creates a HTTP scraper for jobs.
func NewHTTPScraper(config Config) (scrape.JobsScraper, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	tlsConfig, err := createTLSConfig(config)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &httpJobsScraper{
		tlsConfig:            tlsConfig,
		httpClient:           httpClient,
		baseURL:              config.BaseURL,
		runContext:           ctx,
		runContextCancelFunc: cancel,
	}, nil
}

// createTLSConfig creates a TLS config. Should only be called after config.Validate().
func createTLSConfig(config Config) (tlsConfig *tls.Config, err error) {
	if !strings.HasPrefix(config.BaseURL, "https://") {
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
	if config.caCertPool != nil {
		tlsConfig.RootCAs = config.caCertPool
	}
	return tlsConfig, nil
}
