package api

import (
	"github.com/janoszen/openshiftci-inspector/common/api"
	"github.com/janoszen/openshiftci-inspector/jobs"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

func NewJobsListAPI(jobsStorage storage.JobsStorage) api.API {
	return &jobsListAPI{
		storage: jobsStorage,
	}
}

type jobsListAPI struct {
	storage storage.JobsStorage
}

func (j *jobsListAPI) GetRoutes() []api.Route {
	return []api.Route{
		{
			Method: "GET",
			Path:   "/jobs",
		},
	}
}

// Handle returns a list of jobs currently stored.
//
// swagger:route GET /jobs jobs listJobs
//
// Get a list of jobs currently stored.
//
// Produces:
// - application/json
//
// Consumes:
// - application/json
//
// Schemes: http
//
// Responses:
// default: JobsListResponse
//
func (j *jobsListAPI) Handle(_ api.Request, response api.Response) error {
	jobList, err := j.storage.ListJobs()
	if err != nil {
		return err
	}
	return response.Encode(JobsListResponseBody{
		Jobs: jobList,
	})
}

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
