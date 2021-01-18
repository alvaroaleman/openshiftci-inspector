package indexer

import (
	"context"
	"sync"

	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

func NewJobsIndexer(storage storage.JobsStorage) JobIndexer {
	runContext, cancelRunContext := context.WithCancel(context.Background())
	return &jobsIndexer{
		storage:          storage,
		runContext:       runContext,
		cancelRunContext: cancelRunContext,
		done:             make(chan struct{}),
		mu:               &sync.Mutex{},
	}
}
