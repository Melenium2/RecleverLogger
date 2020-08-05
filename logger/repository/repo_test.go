package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RecleverLogger/logger"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var source = os.Getenv("db_url")

func initRealDb(t *testing.T) (*sqlx.DB, func(...string)) {
	db := CreateDatabase(source, "../../")

	return db, func(tables ...string) {
		tr, err := db.Begin()
		assert.NoError(t, err)
		for _, tab := range tables {
			tr.Exec("TRUNCATE TABLE ?", tab)
		}
		assert.NoError(t, tr.Commit())
	}
}

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

func TestLog_MockSave_ShouldSaveNewLogToDb(t *testing.T) {
	ctx := context.Background()
	db, mock, err := initTestDb()
	assert.NoError(t, err)

	now := time.Now().UTC()
	l := &logger.SingleLog{
		Type:       "info",
		Module:     "backend",
		Message:    "hello",
		Stacktrace: "big trace",
		Time:       now,
		Timestamp:  now.Unix(),
	}
	mock.ExpectExec("^insert into logs \\(type, module, message, stacktrace, time, timestamp\\) values \\(\\?, \\?, \\?, \\?, \\?, \\?\\)$").
		WithArgs(l.Type, l.Module, l.Message, l.Stacktrace, l.Time, l.Timestamp).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := New(db, NopLogger{})
	assert.NoError(t, repo.Save(ctx, l))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLog_Save_ShouldSaveNewLogToDb(t *testing.T) {
	if source == "" {
		t.Fatal("For test TestLog_Save_ShouldSaveNewLogToDb need to provide db_url env var")
	}
	db, cleaner := initRealDb(t)
	defer cleaner("logs")
	ctx := context.Background()
	repo := New(db, NopLogger{})

	now := time.Now().UTC()
	l := &logger.SingleLog{
		Type:       "info",
		Module:     "backend",
		Message:    "hello",
		Stacktrace: "big trace",
		Time:       now,
		Timestamp:  now.Unix(),
	}

	assert.NoError(t, repo.Save(ctx, l))
}
