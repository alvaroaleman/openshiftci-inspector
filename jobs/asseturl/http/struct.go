package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/janoszen/openshiftci-inspector/jobs"
)

type httpAssetURLFetcher struct {
	httpClient *http.Client
	exit       chan struct{}
	done       chan struct{}
	baseURL    string
}

func (h *httpAssetURLFetcher) Process(input <-chan jobs.Job) <-chan jobs.JobWithAssetURL {
	artifactsRe := regexp.MustCompile(`<a href="(?P<url>[^"]+)">Artifacts</a>`)
	result := make(chan jobs.JobWithAssetURL)
	go func() {
		defer func() {
			_ = recover()
			close(h.exit)
			close(h.done)
		}()
	loop:
		for {
			select {
			case <-h.exit:
				break loop
			case job, ok := <-input:
				if !ok {
					break loop
				}
				jobPage, err := http.Get(job.URL)
				if err != nil {
					// TODO log and retry
					continue
				}
				body, err := ioutil.ReadAll(jobPage.Body)
				if err != nil {
					//TODO log and retry
					continue
				}
				matches := artifactsRe.FindStringSubmatch(string(body))
				if len(matches) > 1 {
					result <- jobs.JobWithAssetURL{
						Job:      job,
						AssetURL: matches[1],
					}
				}
				// TODO log and retry if no match
			}
		}
	}()
	return result
}

func (h *httpAssetURLFetcher) Shutdown(ctx context.Context) {
	select {
	case <-h.done:
		return
	case <-ctx.Done():
		close(h.exit)
		<-h.done
	}
}
