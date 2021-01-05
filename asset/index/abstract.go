package index

import (
	"context"
)

// AbstractAssetIndex is a default implementation for AssetIndex.
type AbstractAssetIndex struct {
}

// Shutdown is called when the program is shutting down and gives the asset index a chance to clean up.
func (a *AbstractAssetIndex) Shutdown(_ context.Context) {
}
