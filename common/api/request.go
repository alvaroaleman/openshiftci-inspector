package api

import (
	"net/http"
)

type request struct {
	decoders []Decoder
	pathVars map[string]string
	request  *http.Request
}

func (r *request) Decode(target interface{}) error {
	contentType := r.request.Header.Get("content-type")
	bodyDecoded := false
	for _, decoder := range r.decoders {
		bodyDecoder := decoder.BodyDecoder()
		var err error
		if bodyDecoder != nil {
			if contentType != "" && bodyDecoder.CanDecode(contentType) {
				err = bodyDecoder.Decode(r.pathVars, r.request, target)
				bodyDecoded = true
			}
		} else {
			err = decoder.Decode(r.pathVars, r.request, target)
		}
		if err != nil {
			return err
		}
	}
	if r.request.Method != "GET" && !bodyDecoded {
		return &HttpError{
			StatusCode: http.StatusUnsupportedMediaType,
		}
	} else {
		return nil
	}
}
