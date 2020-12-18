package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"k8s.io/api/core/v1"
)

type JobMetadata struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	SelfLink string `json:"selfLink"`
	UID string `json:"uid"`
	ResourceVersion string `json:"resourceVersion"`
	Generation int `json:"generation"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}

type JobSpec struct {
	Type string `json:"type"`
	Agent string `json:"agent"`
	Cluster string `json:"cluster"`
	Namespace string `json:"namespace"`
	Job string `json:"job"`
	Refs JobRefs `json:"refs"`
	Report bool `json:"report"`
	Context string `json:"context"`
	RerunCommand string `json:"rerun_command"`
	PodSpec v1.PodSpec `json:"pod_spec"`
}

type JobRefs struct {
	Org string `json:"org"`
	Repo string `json:"repo"`
	RepoLink string `json:"repo_link"`
	BaseRef string `json:"base_ref"`
	BaseSha string `json:"base_sha"`
	BaseLink string `json:"base_link"`
	Pulls []JobPull `json:"pulls"`
}

type JobStatus struct {
	StartTime *time.Time `json:"startTime,omitempty"`
	PendingTime *time.Time `json:"pendingTime,omitempty"`
	CompletionTime *time.Time `json:"completionTime,omitempty"`
	State string `json:"state"`
	Description string `json:"description"`
	URL string `json:"url"`
	PodName string `json:"pod_name"`
	BuildID string `json:"build_id"`
}

type JobPull struct {
	Number int `json:"number"`
	Author string `json:"author"`
	Sha string `json:"sha"`
	CommitLink string `json:"commit_link"`
	AuthorLink string `json:"author_link"`
}

type Job struct {
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Metadata JobMetadata `json:"metadata"`
	Spec JobSpec `json:"spec"`
	Status JobStatus `json:"status"`
}

type Jobs struct {
	Items []Job `json:"items"`
}

func FetchJobs() ([]Job, error) {
	url := "https://prow.ci.openshift.org/prowjobs.js?var=allBuilds&omit=annotations,labels,decoration_config,pod_spec"
	client := http.Client{}
	data, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	rawJSON := strings.Replace(string(bytes[:len(bytes)-1]), "var allBuilds = ", "", 1)

	jobs := Jobs{}
	if err := json.Unmarshal([]byte(rawJSON), &jobs); err != nil {
		return nil, err
	}
	return jobs.Items, err
}
