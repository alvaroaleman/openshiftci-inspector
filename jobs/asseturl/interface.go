package asseturl

import (
	"github.com/janoszen/openshiftci-inspector/common"
	"github.com/janoszen/openshiftci-inspector/jobs"
)

// JobAssetURLFetcher scrapes the asset URL from the job page if not already stored.
type JobAssetURLFetcher interface {
	common.ShutdownHandler

	// Process receives job records and then scrapes the asset URLs and offers them on the second channel.
	Process(<-chan jobs.Job) <-chan jobs.JobWithAssetURL
}

type JobAssetURLScraper interface {
}
