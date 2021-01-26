package storage

import (
	"errors"
	"time"

	"github.com/janoszen/openshiftci-inspector/common"
	"github.com/janoszen/openshiftci-inspector/jobs"
)

var ErrJobNotFound = errors.New("job not found")

type ListJobsParams struct {
	Job        *string
	GitOrg     *string
	GitRepo    *string
	PullNumber *int
	Before     *time.Time
	After      *time.Time
}

// JobsStorage stores a database of all jobs.
type JobsStorage interface {
	common.ShutdownHandler

	// UpdateJob updates the data on a single job.
	UpdateJob(job jobs.Job) (err error)

	// ListJobs lists all jobs.
	ListJobs(ListJobsParams) ([]jobs.Job, error)

	// GetJob returns a single job.
	GetJob(id string) (jobs.Job, error)
}

// ErrJobHasNoAssetURL is an error that indicates that the specified job has no stored asset URL and the asset URL should
// be fetched from the job URL.
var ErrJobHasNoAssetURL = errors.New("the requested job has no asset URL")

// JobsAssetURLStorage is a storage interface that lets you update and fetch the asset URL.
type JobsAssetURLStorage interface {
	common.ShutdownHandler

	// UpdateAssetURL sets the assetURL to the specified value for a job.
	UpdateAssetURL(job jobs.Job, assetURL string) error

	// GetAssetURLForJob returns the asset URL for a job if present or a ErrJobHasNoAssetURL if not found. It can also
	// return other errors if the fetch failed.
	GetAssetURLForJob(job jobs.Job) (assetURL string, err error)
}

// CompoundJobsStorage combines the JobsStorage and the JobsAssetURLStorage
type CompoundJobsStorage interface {
	JobsStorage
	JobsAssetURLStorage
}
