package downloader

import (
	"github.com/janoszen/openshiftci-inspector/asset"
)

// AssetDownloader is responsible for downloading assets based on the records.
type AssetDownloader interface {
	// Download receives asset records through a channel and downloads the asset.
	Download(asset asset.AssetWithJob) error
}
