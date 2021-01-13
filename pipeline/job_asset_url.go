package pipeline

import (
	"context"

	"github.com/janoszen/openshiftci-inspector/jobs"
)

type JobAssetURLFetcher interface {
	Process(<-chan jobs.Job, chan<- jobs.JobWithAssetURL)
	Shutdown(ctx context.Context)
}
