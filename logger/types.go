package logger

import "time"

type SingleLog struct {
	Type       string    `db:"type" json:"type"`
	Module     string    `db:"module" json:"module,omitempty"`
	Message    string    `db:"message" json:"message,omitempty"`
	Stacktrace string    `db:"stacktrace" json:"stacktrace,omitempty"`
	Time       time.Time `db:"time" json:"time,omitempty"`
	Timestamp  int64     `db:"timestamp" json:"timestamp,omitempty"`
}
