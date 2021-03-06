package scraper

import (
	"context"
	"log"

	"github.com/janoszen/openshiftci_inspector/asset/downloader"
	"github.com/janoszen/openshiftci_inspector/asset/index"
	"github.com/janoszen/openshiftci_inspector/jobs/asseturl"
	"github.com/janoszen/openshiftci_inspector/jobs/scrape"
	"github.com/janoszen/openshiftci_inspector/jobs/storage"
)

func New(
	logger *log.Logger,
	scraper scrape.JobsScraper,
	jobsStorage storage.JobsStorage,
	assetURLFetcher asseturl.JobAssetURLFetcher,
	assetIndex index.AssetIndexer,
	assetDownloader downloader.AssetDownloader,
) Scraper {
	runContext, cancel := context.WithCancel(context.Background())

	return &scraperImpl{
		logger:          logger,
		scraper:         scraper,
		jobsStorage:     jobsStorage,
		assetURLFetcher: assetURLFetcher,
		assetIndex:      assetIndex,
		assetDownloader: assetDownloader,

		runContext:       runContext,
		cancelRunContext: cancel,
	}
}
