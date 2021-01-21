package jobs

import (
	"time"
)

// Job is a description of a CI job.
//
// swagger:model Job
type Job struct {
	// required: true
	ID string `json:"id"`
	// required: true
	Job string `json:"job"`
	// required: true
	JobSafeName string `json:"jobSafeName"`
	// required: true
	Status string `json:"status"`
	// required: false
	StartTime *time.Time `json:"startTime,omitempty"`
	// required: false
	PendingTime *time.Time `json:"pendingTime,omitempty"`
	// required: false
	CompletionTime *time.Time `json:"completionTime,omitempty"`
	// required: true
	URL string `json:"url"`
	// required: false
	GitOrg string `json:"gitOrg,omitempty"`
	// required: false
	GitRepo string `json:"gitRepo,omitempty"`
	// required: false
	GitRepoLink string `json:"gitRepoLink,omitempty"`
	// required: false
	GitBaseRef string `json:"gitBaseRef,omitempty"`
	// required: false
	GitBaseSHA string `json:"gitBaseSHA,omitempty"`
	// required: false
	GitBaseLink string `json:"gitBaseLink,omitempty"`

	// required: true
	Pulls []Pull `json:"pulls"`
}
