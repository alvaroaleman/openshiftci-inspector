package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/janoszen/openshiftci_inspector/frontend"
)

type Server interface {
	Start() error
	Stop(ctx context.Context)
}

func NewServer(
	handlers []API,
	encoders []Encoder,
	decoders []Decoder,
	logger *log.Logger,
) (Server, error) {
	return &server{
		mu:       &sync.Mutex{},
		handlers: handlers,
		logger:   logger,
		decoders: decoders,
		encoders: encoders,
	}, nil
}

type server struct {
	mu       *sync.Mutex
	srv      *http.Server
	handlers []API
	logger   *log.Logger
	decoders []Decoder
	encoders []Encoder
}

func (s *server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		return err
	}
	r := mux.NewRouter()
	subrouters := map[string]*mux.Router{}
	for _, handler := range s.handlers {
		h := handler
		for _, route := range h.GetRoutes() {
			method := route.Method
			path := route.Path
			if _, ok := subrouters[method]; !ok {
				subrouters[method] = r.Methods(route.Method).Subrouter()
			}
			subrouters[method].HandleFunc(
				path,
				func(
					writer http.ResponseWriter,
					request *http.Request,
				) {
					wg := &sync.WaitGroup{}
					wg.Add(1)
					go func() {
						defer wg.Done()
						wrappedRequest := s.wrapRequest(request)
						wrappedResponse := s.wrapResponse(request, writer)
						defer func() {
							p := recover()
							if err, ok := p.(error); ok && err != nil {
								s.handleError(err, wrappedResponse)
							} else if p != nil {
								s.handleError(
									fmt.Errorf("panic while serving request (%v)", p),
									wrappedResponse,
								)
							}
						}()
						err := h.Handle(wrappedRequest, wrappedResponse)
						s.handleError(err, wrappedResponse)
					}()
					wg.Wait()
				})
		}
	}

	fs := frontend.GetFilesystem()
	r.PathPrefix("/").Handler(&customFileServer{
		backend: http.FileServer(fs),
	})

	s.srv = &http.Server{
		Handler: r,
	}
	go func() {
		_ = s.srv.Serve(l)
	}()
	return nil
}

func (s *server) Stop(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.srv != nil {
		_ = s.srv.Shutdown(ctx)
	}
}

func (s *server) wrapRequest(req *http.Request) Request {
	pathVars := mux.Vars(req)
	return &request{
		decoders: s.decoders,
		pathVars: pathVars,
		request:  req,
	}
}

func (s *server) wrapResponse(req *http.Request, writer http.ResponseWriter) Response {
	return &response{
		encoders: s.encoders,
		req:      req,
		writer:   writer,
	}
}

func (s *server) handleError(err error, response Response) {
	if err != nil {
		response.SetStatus(http.StatusInternalServerError)
		s.logger.Printf("Error while handling request (%v)\n", err)
	}
}

type customFileServer struct {
	backend http.Handler
}

type interceptingWriter struct {
	header     http.Header
	statusCode int
	writer     io.Writer
}

func (w *interceptingWriter) Header() http.Header {
	return w.header
}

func (w *interceptingWriter) Write(bytes []byte) (int, error) {
	return w.writer.Write(bytes)
}

func (w *interceptingWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (c *customFileServer) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	buffer := &bytes.Buffer{}
	iWriter := &interceptingWriter{
		header:     map[string][]string{},
		statusCode: 200,
		writer:     buffer,
	}
	c.backend.ServeHTTP(iWriter, r)
	if r.URL.Path != "/" && iWriter.statusCode == 404 {
		r.URL.Path = "/"
		c.backend.ServeHTTP(writer, r)
	} else {
		for headerName, headerValues := range iWriter.header {
			for i, headerValue := range headerValues {
				if i == 0 {
					writer.Header().Set(headerName, headerValue)
				} else {
					writer.Header().Add(headerName, headerValue)
				}
			}
		}
		writer.WriteHeader(iWriter.statusCode)
		_, _ = writer.Write(buffer.Bytes())
	}
}
