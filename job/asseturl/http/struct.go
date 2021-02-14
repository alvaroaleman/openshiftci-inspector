package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/janoszen/openshiftci_inspector/jobs"
)

type httpAssetURLFetcher struct {
	httpClient *http.Client
	baseURL    string
	logger     *log.Logger
}

func (h *httpAssetURLFetcher) Process(job job.Job) (job.JobWithAssetURL, error) {
	artifactsRe := regexp.MustCompile(`<a href="(?P<url>[^"]+)">Artifacts</a>`)

	jobPage, err := http.Get(job.URL)
	if err != nil {
		return job.JobWithAssetURL{}, fmt.Errorf("failed to fetch URL %s (%w)", job.URL, err)
	}
	body, err := ioutil.ReadAll(jobPage.Body)
	if err != nil {
		return job.JobWithAssetURL{}, fmt.Errorf("failed to read from URL %s (%w)", job.URL, err)
	}
	matches := artifactsRe.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		return job.JobWithAssetURL{
			Job:      job,
			AssetURL: matches[1],
		}, nil
	} else {
		return job.JobWithAssetURL{}, fmt.Errorf("no asset URL found for job %s", job.URL)
	}
}
