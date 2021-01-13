package jobs

import (
	"time"
)

// Job is a description of a CI job.
type Job struct {
	ID             string     `json:"id"`
	Job            string     `json:"job"`
	JobSafeName    string     `json:"jobSafeName"`
	Status         string     `json:"status"`
	StartTime      *time.Time `json:"startTime,omitEmpty"`
	PendingTime    *time.Time `json:"pendingTime,omitEmpty"`
	CompletionTime *time.Time `json:"completionTime,omitEmpty"`
	URL            string     `json:"url"`
	GitOrg         string     `json:"gitOrg"`
	GitRepo        string     `json:"gitRepo"`
	GitRepoLink    string     `json:"gitRepoLink"`
	GitBaseRef     string     `json:"gitBaseRef"`
	GitBaseSHA     string     `json:"gitBaseSHA"`
	GitBaseLink    string     `json:"gitBaseLink"`

	Pulls []Pull `json:"pulls"`
}
