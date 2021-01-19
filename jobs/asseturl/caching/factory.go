package caching

import (
	"log"

	"github.com/janoszen/openshiftci-inspector/jobs/asseturl"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

func New(
	backend asseturl.JobAssetURLFetcher,
	storage storage.JobsAssetURLStorage,
	logger *log.Logger,
) asseturl.JobAssetURLFetcher {
	return &cachingJobsAssetURLFetcher{
		storage: storage,
		backend: backend,
		exit:    make(chan struct{}),
		done:    make(chan struct{}),
		logger:  logger,
	}
}
