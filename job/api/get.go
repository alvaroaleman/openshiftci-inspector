package api

import (
	"errors"

	"github.com/janoszen/openshiftci_inspector/common/api"
	"github.com/janoszen/openshiftci_inspector/jobs"
	"github.com/janoszen/openshiftci_inspector/jobs/storage"
)

func NewJobsGetAPI(
	jobsStorage storage.JobsStorage,
	assetURLStorage storage.JobsAssetURLStorage,
) api.API {
	return &jobsGetAPI{
		storage:         jobsStorage,
		assetURLStorage: assetURLStorage,
	}
}

type jobsGetAPI struct {
	storage         storage.JobsStorage
	assetURLStorage storage.JobsAssetURLStorage
}

func (j *jobsGetAPI) GetRoutes() []api.Route {
	return []api.Route{
		{
			Method: "GET",
			Path:   "/jobs/{id:[a-zA-Z0-9-]+}",
		},
	}
}

// Handle returns a job currently stored.
//
// swagger:route GET /jobs/{ID} jobs getJob
//
// Get a single job by ID.
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
// default: SingleJobResponse
//
func (j *jobsGetAPI) Handle(apiRequest api.Request, response api.Response) error {
	request := JobsIDRequest{}
	if err := apiRequest.Decode(&request); err != nil {
		return err
	}
	job, err := j.storage.GetJob(request.ID)
	if err != nil {
		return err
	}
	assetURL, err := j.assetURLStorage.GetAssetURLForJob(job)
	if err != nil && !errors.Is(err, storage.ErrJobHasNoAssetURL) {
		return err
	}
	return response.Encode(
		SingleJobResponseBody{
			Job: job.JobWithAssetURL{
				Job:      job,
				AssetURL: assetURL,
			},
		})
}
