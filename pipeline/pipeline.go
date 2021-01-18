package pipeline

import (
	"github.com/janoszen/openshiftci-inspector/common"
)

type Pipeline interface {
	common.ShutdownHandler

	// Run runs the pipeline until a shutdown signal is received.
	Run()
}
