package storage

import (
	"github.com/janoszen/openshiftci-inspector/asset"
	"github.com/janoszen/openshiftci-inspector/common"
)

// AssetStorage is a simplified API to store and retrieve assets for a job.
// All implementations should embed the AbstractAssetStorage struct in order to facilitate adding new methods later on.
type AssetStorage interface {
	common.ShutdownHandler

	// Store stores an asset in the asset storage and returns an error on failure.
	Store(asset asset.Asset, mime string, data []byte) error

	// Fetch retrieves an asset from the storage and returns it, or an error if the retrieval failed.
	Fetch(asset asset.Asset) ([]byte, error)
}
