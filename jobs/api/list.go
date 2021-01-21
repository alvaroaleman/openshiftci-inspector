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

func (j *jobsListAPI) Handle(_ api.Request, response api.Response) error {
	jobList, err := j.storage.ListJobs()
	if err != nil {
		return err
	}
	return response.Encode(JobsListResponseBody{
		Jobs: jobList,
	})
}

type JobsListResponseBody struct {
	Jobs []jobs.Job
}
