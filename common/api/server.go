package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
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
		for _, route := range handler.GetRoutes() {
			if _, ok := subrouters[route.Method]; !ok {
				subrouters[route.Method] = r.Methods(route.Method).Subrouter()
			}
			subrouters[route.Method].HandleFunc(
				route.Path,
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
						err := handler.Handle(wrappedRequest, wrappedResponse)
						s.handleError(err, wrappedResponse)
					}()
					wg.Wait()
				})
		}
	}

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
