package index

import (
	"github.com/janoszen/openshiftci-inspector/asset/indexstorage"
)

func New(assetIndex indexstorage.AssetIndex) AssetIndexer {
	return &assetIndexer{
		storage: assetIndex,
	}
}
