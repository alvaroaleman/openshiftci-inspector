package index

import (
	"fmt"
	"log"
	"strings"

	"github.com/janoszen/openshiftci_inspector/asset"
	"github.com/janoszen/openshiftci_inspector/asset/indexstorage"
	"github.com/janoszen/openshiftci_inspector/jobs"
)

type assetIndexer struct {
	assets  map[string]map[string]string
	storage indexstorage.AssetIndex
	logger  *log.Logger
}

func (a *assetIndexer) GetMissingAssets(job jobs.JobWithAssetURL) ([]asset.AssetWithJob, error) {
	jobID := job.ID
	var assets []asset.AssetWithJob
	for jobPart, assetList := range a.assets {
		if !strings.Contains(job.Job.Job, jobPart) {
			continue
		}
		for assetName, assetRemotePath := range assetList {
			hasAsset, err := a.storage.HasAsset(jobID, assetName)
			if err != nil {
				return nil, fmt.Errorf(
					"error while checking if asset %s is present for job %s (%w)",
					assetName,
					jobID,
					err,
				)
			}
			if !hasAsset {
				assets = append(
					assets, asset.AssetWithJob{
						Asset: asset.Asset{
							JobID:           jobID,
							AssetName:       assetName,
							AssetRemotePath: assetRemotePath,
						},
						Job: job,
					},
				)
			}
		}
	}
	return assets, nil
}
