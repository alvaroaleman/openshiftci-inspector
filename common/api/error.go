package api

import (
	"net/http"
)

type HttpError struct {
	StatusCode int
}

func (h *HttpError) Error() string {
	return http.StatusText(h.StatusCode)
}
