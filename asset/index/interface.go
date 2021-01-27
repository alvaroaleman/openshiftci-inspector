package index

import (
	"github.com/janoszen/openshiftci-inspector/asset"
	"github.com/janoszen/openshiftci-inspector/jobs"
)

// AssetIndexer is a storage for a list of assets for a job.
type AssetIndexer interface {
	// GetMissingAssets is a processor that checks which assets are present for a job, triggers the
	// retrieval, and then emits a list of assets for a job.
	GetMissingAssets(job jobs.JobWithAssetURL) ([]asset.AssetWithJob, error)
}
