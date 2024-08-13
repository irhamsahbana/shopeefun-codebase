package infrastructure

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitializeLogger will set logging format.
func InitializeLogger(stage string, filename string, logLevel zerolog.Level) {

	var (
		lumberjackLogger = &lumberjack.Logger{
			MaxSize:  100, // megabytes
			MaxAge:   14,  // days
			Filename: filename,
		}
		writers = []io.Writer{zerolog.ConsoleWriter{Out: os.Stderr}, lumberjackLogger}
		mw      = io.MultiWriter(writers...)
	)

	// using json format for production
	var logger zerolog.Logger
	if stage == "production" {
		logger = zerolog.New(lumberjackLogger).With().Timestamp().Caller().Logger().Level(zerolog.InfoLevel)
	} else {
		logger = zerolog.New(mw).With().Timestamp().Caller().Logger().Level(logLevel)
	}
	log.Logger = logger

	q := make(chan os.Signal, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(q, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for {
			<-q
			lumberjackLogger.Close()
			log.Info().Msg("Closing logs ...")
		}
	}()
	go func() {
		for {
			<-c
			if err := lumberjackLogger.Rotate(); err != nil {
				log.Error().Err(err).Msg("Error while rotating logs")
			}
			log.Info().Msg("Rotating logs ...")
		}
	}()
}
