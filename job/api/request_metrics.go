package api

// JobsMetricsRequest is a request to run a query on a job
//
// swagger:parameters getMetrics
type JobsMetricsRequest struct {
	JobsIDRequest

	// In: query
	Query string `query:"query" json:"query"`
}
