package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/janoszen/openshiftci-inspector/jobs"
)

type httpJobsScraper struct {
	httpClient           *http.Client
	baseURL              string
	runContext           context.Context
	runContextCancelFunc func()
}

func (h *httpJobsScraper) Scrape() <-chan jobs.Job {
	jobChannel := make(chan jobs.Job)
	go func() {
		for {
			h.doScrapeRun(jobChannel)
			select {
			case <-time.After(10 * time.Minute):
			case <-h.runContext.Done():
				return
			}
		}
	}()
	return jobChannel
}

func (h *httpJobsScraper) Shutdown(shutdownContext context.Context) {
	h.runContextCancelFunc()
}

func (h *httpJobsScraper) doScrapeRun(jobChannels chan jobs.Job) {
	url := h.baseURL + "/prowjobs.js?var=allBuilds"
	data, err := h.httpClient.Get(url)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(data.Body)
	if err != nil {
		//TODO log
		return
	}
	rawJSON := strings.Replace(string(b[:len(b)-1]), "var allBuilds = ", "", 1)

	jobsRaw := &map[string]interface{}{}
	if err := json.Unmarshal([]byte(rawJSON), jobsRaw); err != nil {
		//TODO log
		return
	}

	list := jobList{}
	if err := json.Unmarshal([]byte(rawJSON), &list); err != nil {
		//TODO log
		return
	}

	for _, rawJob := range list.Items {
		jobNameSafe := ""
		if len(rawJob.Spec.PodSpec.Containers) > 0 {
			for _, env := range rawJob.Spec.PodSpec.Containers[0].Env {
				if env.Name == "JOB_NAME_SAFE" {
					jobNameSafe = env.Value
				}
			}
		}
		job := jobs.Job{
			ID:             rawJob.Metadata.UID,
			Job:            rawJob.Spec.Job,
			JobSafeName:    jobNameSafe,
			Status:         rawJob.Status.State,
			StartTime:      rawJob.Status.StartTime,
			PendingTime:    rawJob.Status.PendingTime,
			CompletionTime: rawJob.Status.CompletionTime,
			URL:            rawJob.Status.URL,
			GitOrg:         rawJob.Spec.Refs.Org,
			GitRepo:        rawJob.Spec.Refs.Repo,
			GitRepoLink:    rawJob.Spec.Refs.RepoLink,
			GitBaseRef:     rawJob.Spec.Refs.BaseRef,
			GitBaseSHA:     rawJob.Spec.Refs.BaseSha,
			GitBaseLink:    rawJob.Spec.Refs.BaseLink,
			Pulls:          []jobs.Pull{},
		}
		for _, p := range rawJob.Spec.Refs.Pulls {
			job.Pulls = append(job.Pulls, jobs.Pull{
				Number:     p.Number,
				Author:     p.Author,
				SHA:        p.SHA,
				PullLink:   p.Link,
				CommitLink: p.CommitLink,
				AuthorLink: p.AuthorLink,
			})
		}
		for {
			select {
			case <-h.runContext.Done():
				return
			case jobChannels <- job:
			}
		}
	}
}
