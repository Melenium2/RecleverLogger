package server

import "github.com/RecleverLogger/server/handlers"

type Config struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	UseTls       bool
	TLSCertFile  string
	TLSKeyFile   string
	Handlers     handlers.Handlers
}