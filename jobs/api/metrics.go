package api

import (
	"context"
	"sync"
	"time"

	assetStorage "github.com/janoszen/openshiftci-inspector/asset/storage"
	"github.com/janoszen/openshiftci-inspector/common/api"
	"github.com/janoszen/openshiftci-inspector/jobs/metrics"
	"github.com/janoszen/openshiftci-inspector/jobs/storage"
)

func NewJobsMetricsAPI(
	jobsStorage storage.JobsStorage,
	assetStorage assetStorage.AssetStorage,
	queryBackend metrics.QueryBackend,
) api.API {
	return &jobMetricsAPI{
		storage:      jobsStorage,
		assetStorage: assetStorage,
		queryBackend: queryBackend,
		lock:         &sync.Mutex{},
	}
}

type jobMetricsAPI struct {
	storage      storage.JobsStorage
	assetStorage assetStorage.AssetStorage
	queryBackend metrics.QueryBackend
	lock         *sync.Mutex
}

func (j *jobMetricsAPI) GetRoutes() []api.Route {
	return []api.Route{
		{
			Method: "GET",
			Path:   "/jobs/{id:[a-zA-Z0-9-]+}/metrics",
		},
	}
}

// Handle returns the queried metrics.
//
// swagger:route GET /jobs/{ID}/metrics jobs getMetrics
//
// Returns returns the queried metrics
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
// default: JobsMetricsResponse
func (j *jobMetricsAPI) Handle(request api.Request, response api.Response) error {
	req := JobsMetricsRequest{}
	if err := request.Decode(&req); err != nil {
		return err
	}

	job, err := j.storage.GetJob(req.ID)
	if err != nil {
		return err
	}

	startTime := time.Now()
	if job.StartTime != nil {
		startTime = *job.StartTime
	}

	endTime := time.Now()
	if job.CompletionTime != nil {
		endTime = *job.CompletionTime
	}

	j.lock.Lock()
	defer j.lock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	queryResult, err := j.queryBackend.Query(ctx, req.ID, "prometheus.tar", req.Query, startTime, endTime)
	if err != nil {
		return err
	}
	result := JobsMetricsResponse{
		Result: queryResult,
	}

	if err := response.Encode(result); err != nil {
		return err
	}
	return nil
}
