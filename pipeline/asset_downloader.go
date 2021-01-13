package pipeline

import (
	"context"

	"github.com/janoszen/openshiftci-inspector/asset"
)

type AssetDownloader interface {
	Download(<-chan asset.AssetWithJob)
	Shutdown(ctx context.Context)
}
