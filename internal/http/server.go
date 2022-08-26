package http

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/hierynomus/iot-monitor/pkg/config"
	"github.com/hierynomus/iot-monitor/pkg/logging"
	"github.com/hierynomus/iot-monitor/pkg/process"
)

var _ process.Process = (*Server)(nil)

type Server struct {
	ctx       context.Context
	config    config.HTTPConfig
	WaitGroup *sync.WaitGroup
	srv       *http.Server
	handlers  map[string]http.Handler
}

func NewServer(ctx context.Context, c config.HTTPConfig) *Server {
	return &Server{
		ctx:       ctx,
		config:    c,
		WaitGroup: &sync.WaitGroup{},
		srv: &http.Server{
			Addr:              c.ListenAddress,
			Handler:           nil,
			ReadTimeout:       c.Timeout,
			ReadHeaderTimeout: c.Timeout,
			WriteTimeout:      c.Timeout,
		},
		handlers: make(map[string]http.Handler),
	}
}

func (s *Server) AddHandler(name string, handler http.Handler) {
	s.handlers[name] = handler
}

func (s *Server) Start(ctx context.Context) error {
	s.WaitGroup.Add(1)

	go s.run(ctx)

	return nil
}

func (s *Server) Stop() error {
	return s.srv.Shutdown(s.ctx)
}

func (s *Server) Wait() {
	s.WaitGroup.Wait()
}

func (s *Server) run(_ context.Context) {
	defer s.WaitGroup.Done()
	logger := logging.LoggerFor(s.ctx, "http-server")
	mux := http.NewServeMux()

	logger.Info().Msg("Starting http server")

	for name, handler := range s.handlers {
		mux.Handle(name, handler)
	}

	s.srv.Handler = mux

	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		logger.Info().Msg("http server stopped")
	} else {
		logger.Error().Err(err).Msg("http server failed")
	}
}
