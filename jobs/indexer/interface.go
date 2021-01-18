package indexer

import (
	"github.com/janoszen/openshiftci-inspector/common"
	"github.com/janoszen/openshiftci-inspector/jobs"
)

// JobIndexer is responsible for storing jobs in a database.
type JobIndexer interface {
	common.ShutdownHandler

	// Index takes jobs from the input, updates the internal database, and then outputs the jobs to the putput.
	Index(input <-chan jobs.Job) (output <-chan jobs.Job)
}
