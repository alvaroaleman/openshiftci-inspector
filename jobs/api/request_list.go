package api

// JobsListRequest is a request for a jobs list.
//
// swagger:parameters listJobs
type JobsListRequest struct {
	// Job name part to search for.
	//
	// In: query
	// required: false
	JobLike string `query:"jobLike" json:"jobLike"`

	// Repository name part to search for.
	//
	// In: query
	// required: false
	RepoLike string `query:"repoLike" json:"repoLike"`

	// How many items to fetch.
	//
	// In: query
	// required: false
	// default: 200
	Limit uint `query:"limit" json:"limit"`

	// At which item to start
	//
	// In: query
	// required: false
	// default: 0
	Offset uint `query:"offset" json:"offset"`
}
