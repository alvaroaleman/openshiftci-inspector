package caching

import (
	"log"

	"github.com/janoszen/openshiftci_inspector/jobs/asseturl"
	"github.com/janoszen/openshiftci_inspector/jobs/storage"
)

type cachingJobsAssetURLFetcher struct {
	storage storage.JobsAssetURLStorage
	backend asseturl.JobAssetURLFetcher
	logger  *log.Logger
}
