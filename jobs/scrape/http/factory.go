package http

import (
	"log"

	http2 "github.com/janoszen/openshiftci-inspector/common/http"
	"github.com/janoszen/openshiftci-inspector/jobs/scrape"
)

// NewHTTPScraper creates a HTTP scraper for jobs.
func NewHTTPScraper(config http2.Config, logger *log.Logger) (scrape.JobsScraper, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	httpClient, err := config.CreateClient()
	if err != nil {
		return nil, err
	}

	return &httpJobsScraper{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		logger:     logger,
	}, nil
}
