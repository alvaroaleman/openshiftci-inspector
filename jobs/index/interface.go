package index

import (
	"context"

	"github.com/janoszen/openshiftci-inspector/jobs"
)

// JobsIndex stores a database of all jobs.
type JobsIndex interface {
	// UpdateJob updates the data on a single job.
	UpdateJob(job jobs.Job)

	// ListJobs lists all jobs.
	ListJobs() ([]jobs.Job, error)

	// Shutdown stops the jobs index.
	Shutdown(shutdownContext context.Context)
}
