package index

import (
	"log"

	"github.com/janoszen/openshiftci-inspector/asset/indexstorage"
)

func New(assetIndex indexstorage.AssetIndex, logger *log.Logger, assets []string) AssetIndexer {
	return &assetIndexer{
		assets:  assets,
		storage: assetIndex,
		logger:  logger,
	}
}
