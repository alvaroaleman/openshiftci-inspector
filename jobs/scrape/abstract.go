package scrape

import (
	"context"
)

// AbstractJobsScraper is the default implementation of the jobs scraper
type AbstractJobsScraper struct {
}

// Shutdown stops the jobs scraper
func (a *AbstractJobsScraper) Shutdown(_ context.Context) {

}
