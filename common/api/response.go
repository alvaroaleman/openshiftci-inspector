package api

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type response struct {
	encoders []Encoder
	writer   http.ResponseWriter
	req      *http.Request
}

func (r *response) SetStatus(code int) {
	r.writer.WriteHeader(code)
}

func (r *response) SetHeader(header string, value string) {
	r.writer.Header().Add(header, value)
}

func (r *response) Encode(response interface{}) error {
	acceptMimes := r.extractSortedAcceptMimeList()
	for _, mime := range acceptMimes {
		for _, encoder := range r.encoders {
			if encoder.SupportsMime(mime) {
				return encoder.Encode(response, r.writer)
			}
		}
	}
	r.SetStatus(http.StatusNotAcceptable)
	return nil
}

func (r *response) extractSortedAcceptMimeList() []string {
	acceptHeader := r.req.Header.Get("accept")
	if acceptHeader == "" {
		acceptHeader = "text/html"
	}
	acceptHeaderParts := strings.Split(acceptHeader, ",")
	var acceptHeaders []struct {
		mime string
		q    float64
	}
	for _, accept := range acceptHeaderParts {
		foundQ := false
		acceptParts := strings.Split(accept, ";")
		if len(acceptParts) > 1 {
			for _, part := range acceptParts[1:] {
				kv := strings.SplitN(part, "=", 2)
				if len(kv) == 2 {
					if kv[0] == "q" {
						q, err := strconv.ParseFloat(kv[1], 64)
						if err == nil {
							foundQ = true
							acceptHeaders = append(
								acceptHeaders, struct {
									mime string
									q    float64
								}{acceptParts[0], q},
							)
						}
					}
				}
			}
		}
		if !foundQ {
			acceptHeaders = append(
				acceptHeaders, struct {
					mime string
					q    float64
				}{acceptParts[0], 1},
			)
		}
	}
	sort.SliceStable(
		acceptHeaders, func(i, j int) bool {
			return acceptHeaders[i].q > acceptHeaders[j].q
		},
	)
	var acceptMimes []string
	for _, s := range acceptHeaders {
		acceptMimes = append(acceptMimes, s.mime)
	}
	return acceptMimes
}
