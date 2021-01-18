package caching

import (
	"context"
	"errors"

	"github.com/janoszen/openshiftci-inspector/jobs"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

func (h *httpJobsAssetURLFetcher) Process(input <-chan jobs.Job) <-chan jobs.JobWithAssetURL {
	// TODO periodically retry and reinject jobs without asset URLs.

	backendQueue := make(chan jobs.Job)
	backendReturn := h.backend.Process(backendQueue)

	result := make(chan jobs.JobWithAssetURL)
	go func() {
		defer func() {
			_ = recover()
			//TODO log panic
			close(backendQueue)
			close(result)
			close(h.done)
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
			case <-h.exit:
				break loop
			}
			assetURL, err := h.storage.GetAssetURLForJob(job)
			if err != nil {
				if errors.Is(err, storage.JobHasNoAssetURL) {
					backendQueue <- job
					jobResult, ok := <-backendReturn
					if !ok {
						break loop
					}
					h.storage.UpdateAssetURL(job, jobResult.AssetURL)
					assetURL = jobResult.AssetURL
				} else {
					continue
				}
			}
			if assetURL != "" {
				result <- jobs.JobWithAssetURL{
					Job:      job,
					AssetURL: assetURL,
				}
			}
		}
	}()
	return result
}

func (h *httpJobsAssetURLFetcher) Shutdown(ctx context.Context) {
	select {
	case <-h.done:
		return
	case <-ctx.Done():
		close(h.exit)
	}
}
