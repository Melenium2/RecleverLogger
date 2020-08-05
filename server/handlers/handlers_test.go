package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RecleverLogger/logger"
	"github.com/RecleverLogger/logger/repository"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func initTestDb() (*sqlx.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	xdb := sqlx.NewDb(db, "sqlmock")
	return xdb, mock, nil
}

type NopLogger struct {}

func (n NopLogger) Logs(log ...interface{}) error {
	return nil
}

func (n NopLogger) Logf(message string, log ...interface{}) error {
	return nil
}

func TestHandlers_Log_ShouldReturnNoError(t *testing.T) {
	now := time.Now().UTC()
	l := &logger.SingleLog{
		Type:       "info",
		Module:     "backend",
		Message:    "hello",
		Stacktrace: "big trace",
		Time:       now,
		Timestamp:  now.Unix(),
	}
	b := &bytes.Buffer{}
	assert.NoError(t, json.NewEncoder(b).Encode(l))
	req := httptest.NewRequest("POST", "/log", b)
	resp := httptest.NewRecorder()

	db, mock, err := initTestDb()
	assert.NoError(t, err)

	repo := repository.New(db, NopLogger{})
	mock.ExpectExec("^insert into logs \\(type, module, message, stacktrace, time, timestamp\\) values \\(\\?, \\?, \\?, \\?, \\?, \\?\\)$").
		WithArgs(l.Type, l.Module, l.Message, l.Stacktrace, l.Time, l.Timestamp).
		WillReturnResult(sqlmock.NewResult(1, 1))

	service := &Service{
		db: repo,
		logger: NopLogger{},
	}
	handler := http.HandlerFunc(service.Log)

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "\"Recorded\"\n", resp.Body.String())
}
