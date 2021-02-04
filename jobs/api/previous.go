package api

import (
	"github.com/janoszen/openshiftci_inspector/common/api"
	"github.com/janoszen/openshiftci_inspector/jobs"
	"github.com/janoszen/openshiftci_inspector/jobs/storage"
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
	request := JobsIDLimitRequest{}
	if err := apiRequest.Decode(&request); err != nil {
		return err
	}
	job, err := j.storage.GetJob(request.ID)
	if err != nil {
		return err
	}

	var jobLike *string
	if request.JobLike != "" {
		jobLike = &request.JobLike
	}
	var repoLike *string
	if request.RepoLike != "" {
		repoLike = &request.RepoLike
	}
	limit := uint(200)
	if request.Limit > 0 {
		limit = request.Limit
	}
	offset := uint(0)
	if request.Offset > 0 {
		offset = request.Offset
	}
	jobList, err := j.storage.ListJobs(storage.ListJobsParams{
		Job:      &job.Job,
		GitOrg:   &job.GitOrg,
		GitRepo:  &job.GitRepo,
		Before:   job.StartTime,
		Limit:    &limit,
		Offset:   &offset,
		JobLike:  jobLike,
		RepoLike: repoLike,
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
