package api

import (
	"github.com/janoszen/openshiftci-inspector/common/api"
	"github.com/janoszen/openshiftci-inspector/jobs"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

func NewJobsGetRelatedAPI(
	jobsStorage storage.JobsStorage,
) api.API {
	return &jobsGetRelatedAPI{
		storage: jobsStorage,
	}
}

type jobsGetRelatedAPI struct {
	storage storage.JobsStorage
}

func (j *jobsGetRelatedAPI) GetRoutes() []api.Route {
	return []api.Route{
		{
			Method: "GET",
			Path:   "/jobs/{id:[a-zA-Z0-9-]+}/related",
		},
	}
}

// Handle returns a list of Related jobs for the same build / branch.
//
// swagger:route GET /jobs/{ID}/related jobs getRelatedJobs
//
// Returns a list of related jobs for the same build and branch.
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
func (j *jobsGetRelatedAPI) Handle(apiRequest api.Request, response api.Response) error {
	request := JobsIDLimitRequest{}
	if err := apiRequest.Decode(&request); err != nil {
		return err
	}
	job, err := j.storage.GetJob(request.ID)
	if err != nil {
		return err
	}

	var pullNumber *int
	if job.Pulls != nil && len(job.Pulls) > 0 {
		pullNumber = &job.Pulls[0].Number
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

	jobList, err := j.storage.ListJobs(
		storage.ListJobsParams{
			GitOrg:     &job.GitOrg,
			GitRepo:    &job.GitRepo,
			PullNumber: pullNumber,
			Limit:      &limit,
			Offset:     &offset,
			JobLike:    jobLike,
			RepoLike:   repoLike,
		},
	)
	if err != nil {
		return err
	}
	if jobList == nil {
		jobList = []jobs.Job{}
	}
	return response.Encode(
		JobsListResponseBody{
			Jobs: jobList,
		},
	)
}
