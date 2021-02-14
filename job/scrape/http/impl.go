package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/janoszen/openshiftci_inspector/jobs"
)

type httpJobsScraper struct {
	httpClient *http.Client
	baseURL    string
	logger     *log.Logger
}

func (h *httpJobsScraper) Scrape() ([]job.Job, error) {
	url := h.baseURL + "/prowjobs.js?var=allBuilds"
	data, err := h.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	rawJSON := strings.Replace(string(b[:len(b)-1]), "var allBuilds = ", "", 1)

	jobsRaw := &map[string]interface{}{}
	if err := json.Unmarshal([]byte(rawJSON), jobsRaw); err != nil {
		return nil, err
	}

	list := jobList{}
	if err := json.Unmarshal([]byte(rawJSON), &list); err != nil {
		return nil, err
	}

	i := 0
	var jobList []job.Job
	for _, rawJob := range list.Items {
		i++
		jobNameSafe := ""
		if len(rawJob.Spec.PodSpec.Containers) > 0 {
			for _, env := range rawJob.Spec.PodSpec.Containers[0].Env {
				if env.Name == "JOB_NAME_SAFE" {
					jobNameSafe = env.Value
				}
			}
		}
		job := job.Job{
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
			Pulls:          []job.Pull{},
		}
		for _, p := range rawJob.Spec.Refs.Pulls {
			job.Pulls = append(job.Pulls, job.Pull{
				Number:     p.Number,
				Author:     p.Author,
				SHA:        p.SHA,
				PullLink:   p.Link,
				CommitLink: p.CommitLink,
				AuthorLink: p.AuthorLink,
			})
		}
		jobList = append(jobList, job)
	}
	return jobList, nil
}
