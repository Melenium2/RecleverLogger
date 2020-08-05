package repository

import (
	"context"
	"fmt"
	"github.com/RecleverLogger/logger"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewLog(typ, module, message, stacktrace string) *logger.SingleLog {
	t := time.Now().UTC()
	return &logger.SingleLog{
		Type: typ,
		Module: module,
		Message: message,
		Stacktrace: stacktrace,
		Time: t,
		Timestamp: t.Unix(),
	}
}

type Logs interface {
	Save(context.Context, *logger.SingleLog) error
}

type LoggerRepository struct {
	logger logger.Logger
	db     *sqlx.DB
}

func New(db *sqlx.DB, log logger.Logger) Logs {
	return &LoggerRepository{
		db:     db,
		logger: log,
	}
}

func (c *LoggerRepository) Save(ctx context.Context, log *logger.SingleLog) error {
	if _, err := c.db.ExecContext(
		ctx,
		fmt.Sprint("insert into logs (type, module, message, stacktrace, time, timestamp) values (?, ?, ?, ?, ?, ?)"),
		log.Type, log.Module, log.Message, log.Stacktrace, log.Time, log.Timestamp,
	); err != nil {
		c.logger.Logs("[Error]", err)
		return err
	}

	return nil
}
