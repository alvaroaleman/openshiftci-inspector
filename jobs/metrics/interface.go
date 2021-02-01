package metrics

import (
	"context"
	"time"
)

type QueryBackend interface {
	Query(
		ctx context.Context,
		jobID string,
		name string,
		query string,
		startTime time.Time,
		endTime time.Time,
	) (QueryResponse, error)
}

// QueryResponse is the complete result of the query.
//
// swagger:model QueryResponse
type QueryResponse struct {
	Vector []QuerySample `json:"vector,omitempty"`
	Matrix []QuerySeries `json:"matrix,omitempty"`
	Scalar QueryPoint    `json:"scalar,omitempty"`
}

// QueryLabel is a single label name and value.
//
// swagger:model QueryLabel
type QueryLabel struct {
	// required: true
	Name string `json:"name"`
	// required: true
	Value string `json:"value"`
}

// QueryPoint is a single timestamp-value pair.
//
// swagger:model QueryPoint
type QueryPoint struct {
	// required: true
	Timestamp int64 `json:"timestamp"`
	// required: true
	Value float64 `json:"value"`
}

// QuerySeries is a single number series with labels.
//
// swagger:model QuerySeries
type QuerySeries struct {
	// required: true
	Labels []QueryLabel `json:"labels"`
	// required: true
	Points []QueryPoint `json:"points"`
}

// QuerySerie is a single number with a label.
//
// swagger:model QuerySample
type QuerySample struct {
	// required: true
	Labels []QueryLabel `json:"labels"`
	// required: true
	Point QueryPoint `json:"point"`
}
