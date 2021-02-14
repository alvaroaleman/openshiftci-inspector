package api

import (
	"github.com/janoszen/openshiftci_inspector/jobs"
)

// SingleJobResponse is the response to a request to get a single job in the Openshift CI.
//
// swagger:response SingleJobResponse
type SingleJobResponse struct {
	// JobsGetResponseBody is the response body.
	//
	// In: body
	// required: true
	JobsGetResponseBody SingleJobResponseBody `json:",inline"`
}

// SingleJobResponseBody represents a response with a single job.
//
// swagger:model SingleJobResponseBody
type SingleJobResponseBody struct {
	// Job is a single job record.
	//
	// required: true
	Job job.JobWithAssetURL `json:"job"`
}
