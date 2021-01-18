package pipeline

import (
	"context"

	"github.com/janoszen/openshiftci-inspector/asset/downloader"
	"github.com/janoszen/openshiftci-inspector/asset/index"
	"github.com/janoszen/openshiftci-inspector/jobs/asseturl"
	"github.com/janoszen/openshiftci-inspector/jobs/indexer"
	"github.com/janoszen/openshiftci-inspector/jobs/scrape"
)

func New(
	scraper scrape.JobsScraper,
	jobIndexer indexer.JobIndexer,
	assetURLFetcher asseturl.JobAssetURLFetcher,
	assetIndex index.AssetIndexer,
	assetDownloader downloader.AssetDownloader,
) Pipeline {
	runContext, cancel := context.WithCancel(context.Background())

	return &pipelineImpl{
		scraper:         scraper,
		jobIndexer:      jobIndexer,
		assetURLFetcher: assetURLFetcher,
		assetIndex:      assetIndex,
		assetDownloader: assetDownloader,

		runContext:       runContext,
		cancelRunContext: cancel,
	}
}
