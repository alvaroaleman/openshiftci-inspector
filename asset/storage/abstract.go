package storage

import (
	"context"
)

// AbstractAssetStorage is an abstract implementation for AssetStorage containing some methods' default implementation.
// AssetStorage implementations should embed this struct.
type AbstractAssetStorage struct {
}

// Shutdown is called when the program is shutting down and gives the asset storage a chance to clean up.
// The default implementation does nothing.
func (a *AbstractAssetStorage) Shutdown(_ context.Context) {
}
