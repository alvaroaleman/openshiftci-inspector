package api

import (
	"github.com/janoszen/openshiftci-inspector/common/api"
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
	jobList, err := j.storage.ListJobs(storage.ListJobsParams{})
	if err != nil {
		return err
	}
	return response.Encode(JobsListResponseBody{
		Jobs: jobList,
	})
}
