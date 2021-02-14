package api

// JobsIDRequest is a request to fetch a single job.
//
// swagger:parameters getJob
type JobsIDRequest struct {
	// ID of the job to fetch.
	//
	// In: path
	// required: true
	ID string `path:"id"`
}

// JobsIDLimitRequest is a request to fetch a list of items within a single job
//
// swagger:parameters getPreviousJobs getRelatedJobs
type JobsIDLimitRequest struct {
	JobsIDRequest
	JobsListRequest
}
