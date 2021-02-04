package caching

import (
	"errors"
	"fmt"

	"github.com/janoszen/openshiftci_inspector/jobs"
	"github.com/janoszen/openshiftci_inspector/jobs/storage"
)

func (c *cachingJobsAssetURLFetcher) Process(job jobs.Job) (jobs.JobWithAssetURL, error) {
	assetURL, err := c.storage.GetAssetURLForJob(job)
	if err != nil {
		if errors.Is(err, storage.ErrJobHasNoAssetURL) {
			jobResult, err := c.backend.Process(job)
			if err != nil {
				return jobs.JobWithAssetURL{}, fmt.Errorf("failed fetch asset URL for job "+job.ID+" (%w)", err)
			}
			if err := c.storage.UpdateAssetURL(job, jobResult.AssetURL); err != nil {
				return jobs.JobWithAssetURL{}, fmt.Errorf("failed store asset URL for job "+job.ID+" (%w)", err)
			} else {
				assetURL = jobResult.AssetURL
			}
		} else {
			return jobs.JobWithAssetURL{}, fmt.Errorf("failed to get asset URL for job %s (%w)", job.ID, err)
		}
	}
	if assetURL != "" {
		return jobs.JobWithAssetURL{
			Job:      job,
			AssetURL: assetURL,
		}, nil
	}
	return jobs.JobWithAssetURL{}, fmt.Errorf("failed to fetch asset URL")
}
