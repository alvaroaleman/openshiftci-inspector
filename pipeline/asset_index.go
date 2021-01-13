package pipeline

import (
	"context"

	"github.com/janoszen/openshiftci-inspector/asset"
	"github.com/janoszen/openshiftci-inspector/jobs"
)

type AssetIndex interface {
	GetMissingAssets(<-chan jobs.JobWithAssetURL, chan<- asset.AssetWithJob)
	Shutdown(ctx context.Context)
}
