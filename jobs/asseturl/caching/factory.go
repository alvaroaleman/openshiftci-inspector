package caching

import (
	"github.com/janoszen/openshiftci-inspector/jobs/asseturl"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

func New(
	backend asseturl.JobAssetURLFetcher,
	storage storage.JobsAssetURLStorage,
) asseturl.JobAssetURLFetcher {
	return &httpJobsAssetURLFetcher{
		storage: storage,
		backend: backend,
		exit:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}
