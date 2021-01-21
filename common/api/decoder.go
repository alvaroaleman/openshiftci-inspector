package api

import (
	"encoding/json"
	"net/http"
)

type Decoder interface {
	Decode(pathVars map[string]string, request *http.Request, target interface{}) error
	// BodyDecoder returns a BodyDecoder if available, nil otherwise.
	BodyDecoder() BodyDecoder
}

type BodyDecoder interface {
	Decoder

	CanDecode(mime string) bool
}

func NewJSONDecoder() BodyDecoder {
	return &jsonDecoder{}
}

type jsonDecoder struct {
}

func (j *jsonDecoder) BodyDecoder() BodyDecoder {
	return j
}

func (j *jsonDecoder) CanDecode(mime string) bool {
	return mime == "application/json" || mime == "text/json"
}

func (j *jsonDecoder) Decode(_ map[string]string, request *http.Request, target interface{}) error {
	decoder := json.NewDecoder(request.Body)
	return decoder.Decode(target)
}
