package handlers

import (
	"errors"
	"fmt"
	"github.com/RecleverLogger/customerrs"
	"github.com/RecleverLogger/logger"
	"github.com/RecleverLogger/logger/repository"
	"net/http"
)

type Config struct {
	DbUrl                string
	DbInitialMigratePath string
}

type Service struct {
	Handlers Handlers
	db       repository.Logs
	logger   logger.Logger
	config   *Config
}

func New(log logger.Logger, config *Config) (*Service, error) {
	log.Logf("Creating service")

	if config == nil {
		return nil, customerrs.ServiceConfigIsNilErr()
	}
	if config.DbUrl == "" {
		return nil, customerrs.ServiceConfigDbUrlIsEmptyErr()
	}

	service := &Service{
		config: config,
		logger: log,
	}

	{
		db := repository.CreateDatabase(service.config.DbUrl, service.config.DbInitialMigratePath)
		service.db = repository.New(db, log)
	}

	service.Handlers = map[string]Handler{
		"Log": {"/log", service.recoveryWrap(service.Log), "POST"},
	}

	service.logger.Logf("Http service created")
	return service, nil
}

func (s *Service) recoveryWrap(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			var e error
			r := recover()
			if r != nil {
				msg := "Handler recover from the panic"
				switch r.(type) {
				case string:
					e = errors.New(fmt.Sprintf("Error: %s, trace: %s", msg, r))
				case error:
					e = errors.New(fmt.Sprintf("Error: %s, trace: %s", msg, r))
				default:
					e = errors.New(fmt.Sprintf("Error: %s", msg))
				}

				writeError(w, http.StatusInternalServerError, e)
			}
		}()

		if handlerFunc != nil {
			handlerFunc(w, r)
		}
	}
}
