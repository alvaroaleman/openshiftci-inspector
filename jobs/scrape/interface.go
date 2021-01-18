package scrape

import (
	"github.com/janoszen/openshiftci-inspector/common"
	"github.com/janoszen/openshiftci-inspector/jobs"
)

// JobsScraper scrapes jobs off the Prow web UI.
type JobsScraper interface {
	common.ShutdownHandler

	// Scrape supplies a stream of jobs from the Prow web UI.
	Scrape() <-chan jobs.Job
}
