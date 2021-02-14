package api

import (
	"github.com/janoszen/openshiftci_inspector/jobs/metrics"
)

// JobsMetricsResponse is a response to a metrics query.
//
// swagger:response JobsMetricsResponse
type JobsMetricsResponse struct {
	// In: body
	Body JobsMetricsResponseBody `json:",inline"`
}

// JobsMetricsResponseBody represents a response with a metrics reply.
//
// swagger:model
type JobsMetricsResponseBody struct {
	// Result is the result data set to a query.
	//
	// required: true
	metrics.QueryResponse `json:",inline"`
}
