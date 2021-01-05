package storage

import (
	"context"
)

// AssetStorage is a simplified API to store and retrieve assets for a job.
// All implementations should embed the AbstractAssetStorage struct in order to facilitate adding new methods later on.
type AssetStorage interface {
	// Store stores an asset in the asset storage and returns an error on failure.
	Store(jobID string, name string, mime string, data []byte) error

	// Fetch retrieves an asset from the storage and returns it, or an error if the retrieval failed.
	Fetch(jobID string, name string) ([]byte, error)

	// Shutdown is called when the application is terminating, giving the AssetStorage a chance to clean up
	// any pending operations. This method should only return when all shutdown processes are complete.
	// This method should respect the passed shutdownContext parameter and return ASAP when the context is done.
	Shutdown(shutdownContext context.Context)
}
