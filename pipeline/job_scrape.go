package pipeline

import (
	"context"

	"github.com/janoszen/openshiftci-inspector/jobs"
)

type JobScraper interface {
	Scrape(chan<- jobs.Job)
	Shutdown(ctx context.Context)
}
