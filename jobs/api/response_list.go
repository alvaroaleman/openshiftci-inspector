package api

import (
	"github.com/janoszen/openshiftci_inspector/jobs"
)

// JobsListResponse is the response to a request to list jobs in the Openshift CI.
//
// swagger:response JobsListResponse
type JobsListResponse struct {
	// In: body
	JobsListResponseBody JobsListResponseBody `json:",inline"`
}

// JobsListResponseBody represents a response with a job list.
//
// swagger:model
type JobsListResponseBody struct {
	// Jobs is the list of jobs in the response
	//
	// required: true
	Jobs []jobs.Job `json:"jobs"`
}
