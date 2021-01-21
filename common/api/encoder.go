package api

import (
	"encoding/json"
	"net/http"
)

type Encoder interface {
	SupportsMime(mime string) bool
	Encode(interface{}, http.ResponseWriter) error
}

func NewJSONEncoder() Encoder {
	return &jsonEncoder{}
}

type jsonEncoder struct {
}

func (j jsonEncoder) SupportsMime(mime string) bool {
	return mime == "application/json" ||
		mime == "application/*" ||
		mime == "text/json" ||
		mime == "text/*" ||
		mime == "*/*"
}

func (j jsonEncoder) Encode(i interface{}, writer http.ResponseWriter) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(i)
}
