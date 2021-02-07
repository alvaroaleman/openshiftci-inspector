package http

import (
	"fmt"
	"github.com/janoszen/openshiftci_inspector/jobs"
	"io/ioutil"
	"net/http"
	"regexp"
)

func (h *httpAssetURLFetcher) Process(job jobs.Job) (jobs.JobWithAssetURL, error) {
	artifactsRe := regexp.MustCompile(`<a href="(?P<url>[^"]+)">Artifacts</a>`)

	jobPage, err := http.Get(job.URL)
	if err != nil {
		return jobs.JobWithAssetURL{}, fmt.Errorf("failed to fetch URL %s (%w)", job.URL, err)
	}
	body, err := ioutil.ReadAll(jobPage.Body)
	if err != nil {
		return jobs.JobWithAssetURL{}, fmt.Errorf("failed to read from URL %s (%w)", job.URL, err)
	}
	matches := artifactsRe.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		return jobs.JobWithAssetURL{
			Job:      job,
			AssetURL: matches[1],
		}, nil
	} else {
		return jobs.JobWithAssetURL{}, fmt.Errorf("no asset URL found for job %s", job.URL)
	}
}
