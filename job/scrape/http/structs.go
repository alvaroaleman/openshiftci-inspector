package http

import (
	"time"

	v1 "k8s.io/api/core/v1"
)

type jobMetadata struct {
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace"`
	SelfLink          string    `json:"selfLink"`
	UID               string    `json:"uid"`
	ResourceVersion   string    `json:"resourceVersion"`
	Generation        int       `json:"generation"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}

type jobSpec struct {
	Type         string     `json:"type"`
	Agent        string     `json:"agent"`
	Cluster      string     `json:"cluster"`
	Namespace    string     `json:"namespace"`
	Job          string     `json:"job"`
	Refs         jobRefs    `json:"refs"`
	Report       bool       `json:"report"`
	Context      string     `json:"context"`
	RerunCommand string     `json:"rerun_command"`
	PodSpec      v1.PodSpec `json:"pod_spec"`
}

type jobRefs struct {
	Org      string    `json:"org"`
	Repo     string    `json:"repo"`
	RepoLink string    `json:"repo_link"`
	BaseRef  string    `json:"base_ref"`
	BaseSha  string    `json:"base_sha"`
	BaseLink string    `json:"base_link"`
	Pulls    []jobPull `json:"pulls"`
}

type jobStatus struct {
	StartTime      *time.Time `json:"startTime,omitempty"`
	PendingTime    *time.Time `json:"pendingTime,omitempty"`
	CompletionTime *time.Time `json:"completionTime,omitempty"`
	State          string     `json:"state"`
	Description    string     `json:"description"`
	URL            string     `json:"url"`
	PodName        string     `json:"pod_name"`
	BuildID        string     `json:"build_id"`
}

type jobPull struct {
	Number     int    `json:"number"`
	Author     string `json:"author"`
	SHA        string `json:"sha"`
	Link       string `json:"link"`
	CommitLink string `json:"commit_link"`
	AuthorLink string `json:"author_link"`
}

type job struct {
	Kind       string      `json:"kind"`
	ApiVersion string      `json:"apiVersion"`
	Metadata   jobMetadata `json:"metadata"`
	Spec       jobSpec     `json:"spec"`
	Status     jobStatus   `json:"status"`
}

type jobList struct {
	Items []job `json:"items"`
}
