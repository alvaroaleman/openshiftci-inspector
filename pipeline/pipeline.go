package pipeline

import (
	"context"
)

type Pipeline interface {
	// Run runs the pipeline until a shutdown signal is received.
	Run()

	// Shutdown stops the pipeline at the next possible time.
	Shutdown(shutdownContext context.Context)
}
