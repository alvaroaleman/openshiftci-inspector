package indexstorage

import (
	"context"
)

// AssetIndex stores a list of assets for a job.
type AssetIndex interface {
	// AddAsset records an asset for a job.
	AddAsset(jobID string, name string) error

	// HasAsset returns true if the job has an asset with the given name.
	HasAsset(jobID string, name string) (bool, error)

	// ListAssets returns a list of assets for a job.
	ListAssets(jobID string) ([]string, error)

	// Shutdown is called when the program is shutting down and gives the asset index a chance to clean up.
	Shutdown(shutdownContext context.Context)
}
