package log

import (
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

func Logger() zerolog.Logger {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	)

	logLevel := os.Getenv("LOG_LEVEL")
	if strings.ToUpper(logLevel) == "DEBUG" {
		logger = logger.Level(zerolog.DebugLevel)
	} else {
		logger = logger.Level(zerolog.InfoLevel)
	}

	return logger.With().Timestamp().Logger()
}
