package asseturl

import (
	"github.com/janoszen/openshiftci-inspector/jobs"
)

// JobAssetURLFetcher scrapes the asset URL from the job page if not already stored.
type JobAssetURLFetcher interface {
	// Process receives job records and then scrapes the asset URLs and returns it.
	Process(job jobs.Job) (jobs.JobWithAssetURL, error)
}
