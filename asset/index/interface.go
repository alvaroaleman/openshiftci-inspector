package index

import (
	"github.com/janoszen/openshiftci-inspector/asset"
	"github.com/janoszen/openshiftci-inspector/common"
	"github.com/janoszen/openshiftci-inspector/jobs"
)

// AssetIndexer is a storage for a list of assets for a job.
type AssetIndexer interface {
	common.ShutdownHandler

	// GetMissingAssets is a processor that checks which assets are present for a job, triggers the
	// retrieval, and then emits a list of assets for a job.
	GetMissingAssets(<-chan jobs.JobWithAssetURL) <-chan asset.AssetWithJob
}
