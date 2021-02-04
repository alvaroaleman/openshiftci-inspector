package api

import (
	"github.com/janoszen/openshiftci_inspector/common/api"
	"github.com/janoszen/openshiftci_inspector/jobs/storage"
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
func (j *jobsListAPI) Handle(req api.Request, response api.Response) error {
	listRequest := JobsListRequest{}
	if err := req.Decode(&listRequest); err != nil {
		return err
	}

	var jobLike *string
	if listRequest.JobLike != "" {
		jobLike = &listRequest.JobLike
	}
	var repoLike *string
	if listRequest.RepoLike != "" {
		repoLike = &listRequest.RepoLike
	}
	limit := uint(200)
	if listRequest.Limit > 0 {
		limit = listRequest.Limit
	}
	offset := uint(0)
	if listRequest.Offset > 0 {
		offset = listRequest.Offset
	}

	jobList, err := j.storage.ListJobs(storage.ListJobsParams{
		Limit:    &limit,
		Offset:   &offset,
		RepoLike: repoLike,
		JobLike:  jobLike,
	})
	if err != nil {
		return err
	}
	return response.Encode(JobsListResponseBody{
		Jobs: jobList,
	})
}
