package log

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

func Logger() zerolog.Logger {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	)

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "DEBUG" {
		logger = logger.Level(zerolog.DebugLevel)
	}

	return logger.With().Timestamp().Logger()
}
