package scraper

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/janoszen/openshiftci_inspector/asset/downloader"
	"github.com/janoszen/openshiftci_inspector/asset/index"
	"github.com/janoszen/openshiftci_inspector/jobs"
	"github.com/janoszen/openshiftci_inspector/jobs/asseturl"
	"github.com/janoszen/openshiftci_inspector/jobs/scrape"
	"github.com/janoszen/openshiftci_inspector/jobs/storage"
)

type scraperImpl struct {
	logger          *log.Logger
	scraper         scrape.JobsScraper
	jobsStorage     storage.JobsStorage
	assetURLFetcher asseturl.JobAssetURLFetcher
	assetIndex      index.AssetIndexer
	assetDownloader downloader.AssetDownloader

	runContext       context.Context
	cancelRunContext func()
}

func (p *scraperImpl) Run() {
loop:
	for {
		p.logger.Printf("Scraping Prow...")
		jobsList, err := p.scraper.Scrape()
		if err != nil {
			p.logger.Printf("Failed to scrape Prow (%v)", err)
		} else {
			totalJobs := len(jobsList)
			var lock sync.Mutex
			var completeJobs int
			// Matches our db connection
			sem := semaphore.NewWeighted(50)
			for _, job := range jobsList {
				j := job

				if err := sem.Acquire(p.runContext, 1); err != nil {
					break loop
				}
				go func() {
					defer sem.Release(1)
					p.indexJob(j)
					lock.Lock()
					defer lock.Unlock()
					completeJobs++
					printCompletionPercentage(completeJobs, totalJobs)
				}()

				select {
				case <-p.runContext.Done():
					break loop
				default:
				}
			}
			sem.Acquire(p.runContext, 50)
		}

		select {
		case <-p.runContext.Done():
			break loop
		default:
			p.logger.Printf("Sleeping for 10 minutes...")
		}

		select {
		case <-p.runContext.Done():
			break loop
		case <-time.After(10 * time.Minute):
		}
	}
}

func printCompletionPercentage(completeJobs int, totalJobs int) {
	complete := int(100 * float32(completeJobs) / float32(totalJobs))
	bars := ""
	char := "█"
	for i := 0; i < complete; i = i + 2 {
		bars += char
	}
	fmt.Printf("\r[%-50s] %3d%% (%d/%d)", bars, complete, completeJobs, totalJobs)
}

func (p *scraperImpl) Shutdown(_ context.Context) {
	p.cancelRunContext()
}

func (p *scraperImpl) indexJob(job jobs.Job) {
	if err := p.jobsStorage.UpdateJob(job); err != nil {
		p.logger.Printf("\nError while updating job (%v)\033[F", err)
		return
	}

	jobWithAsset, err := p.assetURLFetcher.Process(job)
	if err != nil {
		p.logger.Printf("\nError while fetching asset URL (%v)", err)
		return
	}

	assets, err := p.assetIndex.GetMissingAssets(jobWithAsset)
	if err != nil {
		p.logger.Printf("\nError while generating missing assets for job %v (%v)", job.ID, err)
		return
	}

	for _, asset := range assets {
		err := p.assetDownloader.Download(asset)
		if err != nil {
			p.logger.Printf("\nError while downloading asset (%v)", err)
		}
	}
}
