package http

import (
	"net/http"
)

func (c Config) CreateClient() (*http.Client, error) {
	tlsConfig, err := c.CreateTLSConfig()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{
		Transport: transport,
	}, nil
}
