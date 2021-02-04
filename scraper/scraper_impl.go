package scraper

import (
	"context"
	"fmt"
	"log"
	"time"

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
			completeJobs := 0
			for _, job := range jobsList {
				j := job
				completeJobs++

				p.indexJob(j)

				complete := int(100 * float32(completeJobs) / float32(totalJobs))
				bars := ""
				char := "█"
				for i := 0; i < complete; i = i + 2 {
					bars += char
				}
				fmt.Printf("\r[%-50s] %3d%% (%d/%d)", bars, complete, completeJobs, totalJobs)

				select {
				case <-p.runContext.Done():
					break
				default:
				}
			}
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
		p.logger.Printf("\nError while fetching asset URL (%v)\033[F", err)
		return
	}

	assets, err := p.assetIndex.GetMissingAssets(jobWithAsset)
	if err != nil {
		p.logger.Printf("\nError while loading missing assets (%v)\033[F", err)
		return
	}

	for _, asset := range assets {
		err := p.assetDownloader.Download(asset)
		if err != nil {
			p.logger.Printf("\nError while downloading asset (%v)\033[F", err)
		}
	}
}
