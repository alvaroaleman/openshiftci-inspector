package pipeline

import (
	"context"
	"sync"

	"github.com/janoszen/openshiftci-inspector/asset"
	"github.com/janoszen/openshiftci-inspector/common"
	"github.com/janoszen/openshiftci-inspector/jobs"
)

type pipelineImpl struct {
	scraper         JobScraper
	assetDownloader AssetDownloader
	assetURLFetcher JobAssetURLFetcher
	assetIndex      AssetIndex
	jobIndexer      JobIndexer

	runContext       context.Context
	cancelRunContext func()
}

func (p *pipelineImpl) Run() {
	jobsChannel := make(chan jobs.Job)
	indexedJobsChannel := make(chan jobs.Job)
	jobsWithAssetURLChannel := make(chan jobs.JobWithAssetURL)
	assetWithJobPipeline := make(chan asset.AssetWithJob)

	go p.scraper.Scrape(jobsChannel)
	go p.jobIndexer.Index(jobsChannel, indexedJobsChannel)
	go p.assetURLFetcher.Process(indexedJobsChannel, jobsWithAssetURLChannel)
	go p.assetIndex.GetMissingAssets(jobsWithAssetURLChannel, assetWithJobPipeline)
	go p.assetDownloader.Download(assetWithJobPipeline)

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
