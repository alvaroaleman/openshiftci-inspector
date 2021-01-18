package pipeline

import (
	"context"
	"sync"

	"github.com/janoszen/openshiftci-inspector/asset/downloader"
	"github.com/janoszen/openshiftci-inspector/asset/index"
	"github.com/janoszen/openshiftci-inspector/common"
	"github.com/janoszen/openshiftci-inspector/jobs/asseturl"
	"github.com/janoszen/openshiftci-inspector/jobs/indexer"
	"github.com/janoszen/openshiftci-inspector/jobs/scrape"
)

type pipelineImpl struct {
	scraper         scrape.JobsScraper
	jobIndexer      indexer.JobIndexer
	assetURLFetcher asseturl.JobAssetURLFetcher
	assetIndex      index.AssetIndexer
	assetDownloader downloader.AssetDownloader

	runContext       context.Context
	cancelRunContext func()
}

func (p *pipelineImpl) Run() {
	jobsChannel := p.scraper.Scrape()
	indexedJobsChannel := p.jobIndexer.Index(jobsChannel)
	jobsWithAssetURLChannel := p.assetURLFetcher.Process(indexedJobsChannel)
	assetWithJobPipeline := p.assetIndex.GetMissingAssets(jobsWithAssetURLChannel)
	p.assetDownloader.Download(assetWithJobPipeline)

	<-p.runContext.Done()
}

func (p *pipelineImpl) Shutdown(ctx context.Context) {
	handlers := []common.ShutdownHandler{
		p.scraper,
		p.assetDownloader,
		p.assetURLFetcher,
		p.assetIndex,
		p.jobIndexer,
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(handlers))
	for _, handler := range handlers {
		h := handler
		go func() {
			defer wg.Done()
			h.Shutdown(ctx)
		}()
	}
	wg.Wait()
	p.cancelRunContext()
}
