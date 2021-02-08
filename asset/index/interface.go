package index

import (
	"github.com/janoszen/openshiftci_inspector/asset"
	"github.com/janoszen/openshiftci_inspector/jobs"
)

// AssetIndexer is a storage for a list of assets for a job.
type AssetIndexer interface {
	// GetMissingAssets returns a list of assets that are required but missing for the given job
	GetMissingAssets(job jobs.JobWithAssetURL) ([]asset.AssetWithJob, error)
}
