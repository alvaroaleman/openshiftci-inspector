package pipeline

import (
	"context"

	"github.com/janoszen/openshiftci-inspector/jobs"
)

type JobIndexer interface {
	Index(<-chan jobs.Job, chan<- jobs.Job)
	Shutdown(ctx context.Context)
}
