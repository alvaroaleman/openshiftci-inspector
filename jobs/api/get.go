package api

import (
	"errors"

	"github.com/janoszen/openshiftci-inspector/common/api"
	"github.com/janoszen/openshiftci-inspector/jobs"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
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
// default: JobsGetResponse
//
func (j *jobsGetAPI) Handle(apiRequest api.Request, response api.Response) error {
	request := JobsGetRequest{}
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
	return response.Encode(JobsGetResponseBody{
		Job: jobs.JobWithAssetURL{
			Job:      job,
			AssetURL: assetURL,
		},
	})
}

// JobsGetRequest is a request to fetch a single job.
//
// swagger:parameters getJob
type JobsGetRequest struct {
	// ID of the job to fetch.
	//
	// In: path
	// required: true
	ID string `path:"id"`
}

// JobsGetResponse is the response to a request to get a single job in the Openshift CI.
//
// swagger:response JobsGetResponse
type JobsGetResponse struct {
	// JobsGetResponseBody is the response body.
	//
	// In: body
	// required: true
	JobsGetResponseBody JobsGetResponseBody `json:",inline"`
}

// JobsGetResponseBody represents a response with a single job.
//
// swagger:model JobsGetResponseBody
type JobsGetResponseBody struct {
	// Job is a single job record.
	//
	// required: true
	Job jobs.JobWithAssetURL `json:"job"`
}
