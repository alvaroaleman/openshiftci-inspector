package http

import (
	"log"
	"net/http"
)

type httpAssetURLFetcher struct {
	httpClient *http.Client
	baseURL    string
	logger     *log.Logger
}
