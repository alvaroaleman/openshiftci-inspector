package caching

import (
	"github.com/janoszen/openshiftci-inspector/jobs/asseturl"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

type httpJobsAssetURLFetcher struct {
	storage storage.JobsAssetURLStorage
	backend asseturl.JobAssetURLFetcher
	exit    chan struct{}
	done    chan struct{}
}
