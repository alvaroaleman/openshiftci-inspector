package api

// JobsIDRequest is a request to fetch a single job.
//
// swagger:parameters getJob getPreviousJobs getRelatedJobs
type JobsIDRequest struct {
	// ID of the job to fetch.
	//
	// In: path
	// required: true
	ID string `path:"id"`
}
