package common

import (
	"context"
)

// ShutdownHandler exposes a hook to trigger the stop processing.
type ShutdownHandler interface {
	// Shutdown stops the scraping in progress and shuts down the scrape. The shutdown context gives the deadline by
	// which to shut down. It does not mean a hard stop as each process may have last items to clean up.
	Shutdown(ctx context.Context)
}
