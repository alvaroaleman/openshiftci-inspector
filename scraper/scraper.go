package scraper

import (
	"github.com/janoszen/openshiftci_inspector/common"
)

type Scraper interface {
	common.ShutdownHandler

	// Run runs the pipeline until a shutdown signal is received.
	Run()
}
