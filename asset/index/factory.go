package index

import (
	"github.com/janoszen/openshiftci-inspector/asset/indexstorage"
)

func New(assetIndex indexstorage.AssetIndex, assets []string) AssetIndexer {
	return &assetIndexer{
		assets:  assets,
		storage: assetIndex,
	}
}
