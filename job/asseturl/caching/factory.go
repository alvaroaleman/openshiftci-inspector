package caching

import (
	"log"

	"github.com/janoszen/openshiftci_inspector/jobs/asseturl"
	"github.com/janoszen/openshiftci_inspector/jobs/storage"
)

func New(
	backend asseturl.JobAssetURLFetcher,
	storage storage.JobsAssetURLStorage,
	logger *log.Logger,
) asseturl.JobAssetURLFetcher {
	return &cachingJobsAssetURLFetcher{
		storage: storage,
		backend: backend,
		logger:  logger,
	}
}
