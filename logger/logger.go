package logger

import (
	"fmt"
	murlog "github.com/Melenium2/Murlog"
	"github.com/RecleverLogger/logger/externallogger"
	"time"
)

type Logger interface {
	Logs(log ...interface{}) error
	Logf(message string, log ...interface{}) error
}

type SuperLogger struct {
	logger   murlog.Logger
	external externallogger.ExternalLogger
}

func New(log murlog.Logger, external externallogger.ExternalLogger) Logger {
	return &SuperLogger{
		logger: log,
		external: external,
	}
}

func (s SuperLogger) Logs(log ...interface{}) error {
	if len(log) > 0 {
		defer s.external.Sendlog(0, fmt.Sprintf("%s\n\n%v", timestamp(), format(log[0])))
		s.Logf("%v", log[0])
	}

	return nil
}

func (s SuperLogger) Logf(message string, log ...interface{}) error {
	var msg string
	if len(log) > 0 {
		msg = fmt.Sprintf(message, log...)
	} else {
		msg = message
	}
	s.logger.Log("msg", msg)

	return nil
}

func timestamp() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

func format(log interface{}) string {
	msg := ""
	switch log.(type) {
	case *SingleLog:
		l := log.(*SingleLog)
		msg = fmt.Sprintf("type: %s\nmodule: %s\nmessage: %s\nstacktrace: %s\n", l.Type, l.Module, l.Message, l.Stacktrace)
	default:
		msg = fmt.Sprintf("%v", log)
	}

	return msg
}

