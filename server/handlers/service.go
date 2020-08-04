package handlers

import (
	"errors"
	"fmt"
	"github.com/RecleverLogger/customerrs"
	"github.com/RecleverLogger/logger"
	"github.com/RecleverLogger/logger/repository"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mailru/go-clickhouse"
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
		db := createDatabase(service.config.DbUrl, service.config.DbInitialMigratePath)
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

func createDatabase(dbURL, configDir string) *sqlx.DB {
	log.Print("Connect to db url ", dbURL, " ...")
	c, err := sqlx.Connect("clickhouse", dbURL)
	if err != nil {
		log.Print(err)
		time.Sleep(time.Second * 15)
		createDatabase(dbURL, configDir)
	}
	log.Print("Connected to db.", " Init schema...")

	if configDir != "" {
		ddl, err := ioutil.ReadFile(fmt.Sprintf("%s/config/schema.sql", configDir))
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Read schema from file...")
		if _, err := c.Exec(string(ddl)); err != nil {
			if strings.Contains(err.Error(), "Code: 57") {
				newddl := strings.ReplaceAll(string(ddl), "create table if not exists", "ATTACH TABLE")
				if _, err := c.Exec(newddl); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
		log.Print("Schema created.")
	}

	return c
}
