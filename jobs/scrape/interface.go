package scrape

import (
	"context"

	"github.com/janoszen/openshiftci-inspector/jobs"
)

// JobsScraper scrapes jobs off the Prow web UI.
type JobsScraper interface {
	// Scrape supplies a stream of jobs from the Prow web UI.
	Scrape() <-chan jobs.Job
	// Shutdown stops the scraping in progress and shuts down the scrape.
	Shutdown(shutdownContext context.Context)
}
