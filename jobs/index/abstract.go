package index

import (
	"context"
)

// AbstractJobsIndex is a default implementation for the JobsIndex
type AbstractJobsIndex struct {
}

// Shutdown stops the jobs index.
func (a *AbstractJobsIndex) Shutdown(_ context.Context) {

}

// AbstractJobsAssetURLStorage is a default implementation for JobsAssetURLStorage.
type AbstractJobsAssetURLStorage struct {
}

// Shutdown stops the jobs asset URL storage.
func (a *AbstractJobsAssetURLStorage) Shutdown(_ context.Context) {

}
