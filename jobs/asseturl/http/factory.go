package http

import (
	"github.com/janoszen/openshiftci-inspector/common/http"
	"github.com/janoszen/openshiftci-inspector/jobs/asseturl"
)

func NewHTTPAssetURLFetcher(config http.Config) (asseturl.JobAssetURLFetcher, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	httpClient, err := config.CreateClient()
	if err != nil {
		return nil, err
	}

	return &httpAssetURLFetcher{
		baseURL:    config.BaseURL,
		httpClient: httpClient,
		exit:       make(chan struct{}),
		done:       make(chan struct{}),
	}, nil
}
