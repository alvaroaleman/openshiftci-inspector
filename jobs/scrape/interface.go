package scrape

import (
	"github.com/janoszen/openshiftci_inspector/jobs"
)

// JobsScraper scrapes jobs off the Prow web UI.
type JobsScraper interface {
	// Scrape supplies a list of jobs from the Prow web UI.
	Scrape() ([]jobs.Job, error)
}
