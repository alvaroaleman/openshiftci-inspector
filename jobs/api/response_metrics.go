package api

import (
	"github.com/janoszen/openshiftci-inspector/jobs/metrics"
)

type JobsMetricsResponse struct {
	Result metrics.QueryResponse
}
