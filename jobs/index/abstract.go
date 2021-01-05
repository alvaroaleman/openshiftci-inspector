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
