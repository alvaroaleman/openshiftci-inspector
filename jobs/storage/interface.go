package storage

import (
	"errors"

	"github.com/janoszen/openshiftci-inspector/common"
	"github.com/janoszen/openshiftci-inspector/jobs"
)

// JobsStorage stores a database of all jobs.
type JobsStorage interface {
	common.ShutdownHandler

	// UpdateJob updates the data on a single job.
	UpdateJob(job jobs.Job) (err error)

	// ListJobs lists all jobs.
	ListJobs() ([]jobs.Job, error)
}

// JobHasNoAssetURL is an error that indicates that the specified job has no stored asset URL and the asset URL should
// be fetched from the job URL.
var JobHasNoAssetURL = errors.New("the requested job has no asset URL")

// JobsAssetURLStorage is a storage interface that lets you update and fetch the asset URL.
type JobsAssetURLStorage interface {
	common.ShutdownHandler

	// UpdateAssetURL sets the assetURL to the specified value for a job.
	UpdateAssetURL(job jobs.Job, assetURL string) error

	// GetAssetURLForJob returns the asset URL for a job if present or a JobHasNoAssetURL if not found. It can also
	// return other errors if the fetch failed.
	GetAssetURLForJob(job jobs.Job) (assetURL string, err error)
}

// CompoundJobsStorage combines the JobsStorage and the JobsAssetURLStorage
type CompoundJobsStorage interface {
	JobsStorage
	JobsAssetURLStorage
}
