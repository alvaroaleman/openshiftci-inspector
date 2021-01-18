package downloader

import (
	"github.com/janoszen/openshiftci-inspector/asset"
	"github.com/janoszen/openshiftci-inspector/common"
)

// AssetDownloader is responsible for downloading assets based on the records.
type AssetDownloader interface {
	common.ShutdownHandler

	// Download receives asset records through a channel and downloads the asset.
	Download(<-chan asset.AssetWithJob)
}
