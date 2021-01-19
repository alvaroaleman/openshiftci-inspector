package caching

import (
	"context"
	"errors"

	"github.com/janoszen/openshiftci-inspector/jobs"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

func (c *cachingJobsAssetURLFetcher) Process(input <-chan jobs.Job) <-chan jobs.JobWithAssetURL {
	// TODO periodically retry and reinject jobs without asset URLs.

	backendQueue := make(chan jobs.Job)
	backendReturn := c.backend.Process(backendQueue)

	result := make(chan jobs.JobWithAssetURL)
	go func() {
		defer func() {
			_ = recover()
			//TODO log panic
			close(backendQueue)
			close(result)
			close(c.done)
		}()
	loop:
		for {
			var job jobs.Job
			var ok bool
			select {
			case job, ok = <-input:
				if !ok {
					break loop
				}
			case <-c.exit:
				break loop
			}
			c.logger.Println("Fetching asset URL for job " + job.ID + "...")
			assetURL, err := c.storage.GetAssetURLForJob(job)
			if err != nil {
				if errors.Is(err, storage.JobHasNoAssetURL) {
					backendQueue <- job
					jobResult, ok := <-backendReturn
					if !ok {
						break loop
					}
					if err := c.storage.UpdateAssetURL(job, jobResult.AssetURL); err != nil {
						c.logger.Println("Failed store asset URL for job " + job.ID + ".")
					} else {
						assetURL = jobResult.AssetURL
					}
				} else {
					c.logger.Printf("Failed to get asset URL for job %s (%v).\n", job.ID, err)
					continue
				}
			}
			if assetURL != "" {
				c.logger.Println("Forwarding job " + job.ID + " with asset URL...")
				result <- jobs.JobWithAssetURL{
					Job:      job,
					AssetURL: assetURL,
				}
			}
		}
	}()
	return result
}

func (h *cachingJobsAssetURLFetcher) Shutdown(ctx context.Context) {
	select {
	case <-h.done:
		return
	case <-ctx.Done():
		close(h.exit)
	}
}
