package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/RecleverLogger/customerrs"
	"github.com/RecleverLogger/logger"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *Config
	errCh  chan<- error
	stopCh chan struct{}

	listner    net.Listener
	router     *mux.Router
	httpServer *http.Server
	logger     logger.Logger
}

func New(ctx context.Context, errCh chan<- error, config *Config, logger logger.Logger) (*Server, error) {
	var err error

	logger.Logf("Creating new server")

	if config == nil {
		return nil, customerrs.ServerConfigIsNilErr()
	}
	if config.Port == "" {
		return nil, customerrs.ServerHttpPortIsEmptyErr()
	}

	server := &Server{
		config: config,
		errCh:  errCh,
		stopCh: make(chan struct{}, 1),
	}

	if ctx == nil {
		server.ctx, server.cancel = context.WithCancel(context.Background())
	} else {
		server.ctx, server.cancel = context.WithCancel(ctx)
	}

	{
		server.httpServer = &http.Server{
			ReadTimeout:  time.Duration(config.ReadTimeout * int(time.Second)),
			WriteTimeout: time.Duration(config.WriteTimeout * int(time.Second)),
			IdleTimeout:  time.Duration(config.WriteTimeout * int(time.Second)),
		}
		// Тут бы tls но пофин и так сойдет
	}
	{
		server.listner, err = net.Listen("tcp", config.Port)
		if err != nil {
			return nil, customerrs.ServerFailToListenPortErr(config.Port, err)
		}
		logger.Logf("Created new listener on port = ", config.Port)
	}
	{
		server.router = mux.NewRouter()
		http.Handle("/", server.router)

		if config.Handlers == nil {
			return nil, customerrs.ServerHaveNoHandlersErr()
		}

		for _, h := range config.Handlers {
			server.router.HandleFunc(h.Path, h.HandleFunc).Methods(h.Method)
			logger.Logf("Register new endpoint path = %s, method = %s", h.Path, h.Method)
		}
	}

	logger.Logf("Http server is created")

	return server, nil
}

func (s *Server) Run() (e error) {

	defer func() {
		r := recover()
		if r != nil {
			msg := "Server recover from panic"
			switch r.(type) {
			case string:
				e = errors.New(fmt.Sprintf("Error: %s, trace: %s", msg, r))
			case error:
				e = errors.New(fmt.Sprintf("Error: %s, trace: %s", msg, r))
			default:
				e = errors.New(fmt.Sprintf("Error: %s", msg))
			}

			s.errCh <- e
		}
	}()

	if s.config.UseTls {
		// Тут стартануть tls серве
	}
	s.logger.Logf("Starting http server")
	return s.httpServer.Serve(s.listner)
}

func (s *Server) Shutdown() (e error) {
	s.logger.Logf("Shutdown the server")
	defer func() { s.stopCh <- struct{}{} }()
	defer s.cancel()

	cancelCtx, cancel := context.WithTimeout(s.ctx, time.Second*30)
	defer cancel()
	if err := s.httpServer.Shutdown(cancelCtx); err != nil {
		e = customerrs.ServerFailedToShutdownErr()
		return
	}

	s.logger.Logf("Server shutdown successfully")
	return
}
