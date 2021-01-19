package index

import (
	"context"
	"log"

	"github.com/janoszen/openshiftci-inspector/asset"
	"github.com/janoszen/openshiftci-inspector/asset/indexstorage"
	"github.com/janoszen/openshiftci-inspector/jobs"
)

type assetIndexer struct {
	assets  []string
	storage indexstorage.AssetIndex
	logger  *log.Logger

	exit chan struct{}
	done chan struct{}
}

func (a *assetIndexer) GetMissingAssets(urls <-chan jobs.JobWithAssetURL) <-chan asset.AssetWithJob {
	result := make(chan asset.AssetWithJob)
	go func() {
		defer func() {
			_ = recover()
			// TODO handle panic
			close(a.exit)
			close(a.done)
		}()

		for {
			select {
			case <-a.exit:
				return
			case url := <-urls:
				jobID := url.ID
				for _, assetName := range a.assets {
					hasAsset, err := a.storage.HasAsset(jobID, assetName)
					if err != nil {
						a.logger.Printf(
							"Error while checking if asset %s is present for job %s (%v)",
							assetName,
							jobID,
							err,
						)
						continue
					}
					if !hasAsset {
						a.logger.Printf(
							"Queueing asset %s for job %s for download...",
							assetName,
							jobID,
						)
						result <- asset.AssetWithJob{
							Asset: asset.Asset{
								JobID:     jobID,
								AssetName: assetName,
							},
							Job: url,
						}
					} else {
						a.logger.Printf(
							"Job %s already has asset %s.",
							jobID,
							assetName,
						)
					}
				}
			}

		}
	}()
	return result
}

func (a *assetIndexer) Shutdown(ctx context.Context) {
	select {
	case <-a.done:
		return
	case <-ctx.Done():
		if a.exit != nil {
			close(a.exit)
		}
		<-a.done
	}
}
