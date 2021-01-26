package api

import (
	"github.com/janoszen/openshiftci-inspector/common/api"
	"github.com/janoszen/openshiftci-inspector/jobs"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

func NewJobsGetPreviousAPI(
	jobsStorage storage.JobsStorage,
) api.API {
	return &jobsGetPreviousAPI{
		storage: jobsStorage,
	}
}

type jobsGetPreviousAPI struct {
	storage storage.JobsStorage
}

func (j *jobsGetPreviousAPI) GetRoutes() []api.Route {
	return []api.Route{
		{
			Method: "GET",
			Path:   "/jobs/{id:[a-zA-Z0-9-]+}/previous",
		},
	}
}

// Handle returns a list of previous jobs for the same build / branch.
//
// swagger:route GET /jobs/{ID}/previous jobs getPreviousJobs
//
// Returns a list of previous jobs for the same build and branch.
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
func (j *jobsGetPreviousAPI) Handle(apiRequest api.Request, response api.Response) error {
	request := JobsIDRequest{}
	if err := apiRequest.Decode(&request); err != nil {
		return err
	}
	job, err := j.storage.GetJob(request.ID)
	if err != nil {
		return err
	}

	jobList, err := j.storage.ListJobs(storage.ListJobsParams{
		Job:     &job.Job,
		GitOrg:  &job.GitOrg,
		GitRepo: &job.GitRepo,
		Before:  job.StartTime,
	})
	if err != nil {
		return err
	}
	if jobList == nil {
		jobList = []jobs.Job{}
	}
	return response.Encode(JobsListResponseBody{
		Jobs: jobList,
	})
}
