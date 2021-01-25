package api

// Route describes a HTTP method and path routed to an API.
//
//swagger:ignore
type Route struct {
	Method string
	Path   string
}

// Request is an interface that handles HTTP requests.
type Request interface {
	// Decode tries to decode the request into the target pointer struct.
	Decode(target interface{}) error
}

// Response is an interface for handling HTTP responses.
type Response interface {
	// SetStatus sets the HTTP status code.
	SetStatus(code int)
	// SetHeader sets a HTTP header in the response.
	SetHeader(header string, value string)
	// Encode takes a response object and encodes if for the client in the requested format.
	Encode(response interface{}) error
}

type API interface {
	GetRoutes() []Route
	Handle(request Request, response Response) error
}
