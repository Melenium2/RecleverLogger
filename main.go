package main

import (
	"context"
	murlog "github.com/Melenium2/Murlog"
	"github.com/RecleverLogger/logger"
	"github.com/RecleverLogger/logger/externallogger"
	"github.com/RecleverLogger/server"
	"github.com/RecleverLogger/server/handlers"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	var (
		httpAddr     = os.Getenv("http_port")
		tgToken      = os.Getenv("tg_token")
		tgChatId     = os.Getenv("tg_chat_id")
		loggerSource = os.Getenv("logger_db")
		configDir    = os.Getenv("config_dir")
	)
	ctx := context.Background()
	errCh := make(chan error, 1)
	sysInter := make(chan os.Signal)
	signal.Notify(sysInter, syscall.SIGINT, syscall.SIGTERM)

	var logger logger.Logger
	var tg *externallogger.TelegramLogger
	{
		id, _ := strconv.Atoi(tgChatId)
		tg = externallogger.NewTelegramLogger(tgToken, id)
		logger = createLogger(tg)
	}

	var handlerConfig *handlers.Config
	var serverConfig *server.Config
	{
		handlerConfig = &handlers.Config{
			DbUrl: loggerSource,
			DbInitialMigratePath: configDir,
		}
		h, err := handlers.New(logger, handlerConfig)
		if err != nil {
			log.Fatal(err)
		}

		serverConfig = &server.Config{
			Port: httpAddr,
			ReadTimeout: 5,
			WriteTimeout: 10,
			IdleTimeout: 30,
			Handlers: h.Handlers,
		}
	}

	s, err := server.New(ctx, errCh, serverConfig, logger)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		errCh <- s.Run()
	}()
	
	select {
	case <-sysInter:
		logger.Logf("System interrupt")
		s.Shutdown()
		return
	case e := <-errCh:
		logger.Logf("Got err %s. Exiting", e)
		s.Shutdown()
		return
	}
}

func createLogger(telegramBot externallogger.ExternalLogger) logger.Logger {
	var l logger.Logger
	{
		var defaultLogger murlog.Logger
		{
			c := murlog.NewConfig()
			c.TimePref(time.RFC1123)
			c.CallerCustomPref(5)
			c.Pref(func() interface{} {
				return "service = globallogger"
			})
			defaultLogger = murlog.NewLogger(c)
		}
		defaultLogger.Log("msg", "Logger db initialized")
		defaultLogger.Log("msg", "Internal logger initialized")

		l = logger.NewLogger(defaultLogger, telegramBot)
	}
	l.Logf("Logger created")
	return l
}