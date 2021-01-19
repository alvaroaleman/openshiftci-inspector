package http

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/url"
	"runtime"
	"strings"
)

// Config is the configuration structure for the HTTP scraper.
type Config struct {
	// BaseURL is the starting URL of the process.
	BaseURL string `json:"baseurl" yaml:"baseurl" env:"HTTP_BASEURL"`
	// CACert is the CA certificate to check against when accessing Prow.
	CACert string `json:"cacert" yaml:"cacert" env:"HTTP_CACERT_FILE"`

	// caCertPool is an internal structure populated by the Validate method.
	caCertPool *x509.CertPool
}

// Validate checks the configuration and loads the certificate in the background.
func (c *Config) Validate() error {
	if _, err := url.ParseRequestURI(c.BaseURL); err != nil {
		return fmt.Errorf("invalid base URL: %s (%w)", c.BaseURL, err)
	}
	if strings.HasPrefix(c.BaseURL, "https://") {
		if strings.TrimSpace(c.CACert) != "" {
			caCert, err := loadPem(c.CACert)
			if err != nil {
				return fmt.Errorf("failed to load CA certificate (%w)", err)
			}

			c.caCertPool = x509.NewCertPool()
			if !c.caCertPool.AppendCertsFromPEM(caCert) {
				return fmt.Errorf("invalid CA certificate provided")
			}
		} else if runtime.GOOS != "windows" {
			//Remove condition if https://github.com/golang/go/issues/16736 gets fixed
			var err error
			c.caCertPool, err = x509.SystemCertPool()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func loadPem(spec string) ([]byte, error) {
	if !strings.HasPrefix(strings.TrimSpace(spec), "-----") {
		return ioutil.ReadFile(spec)
	}
	return []byte(spec), nil
}
