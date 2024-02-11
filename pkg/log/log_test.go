package log

import (
	"github.com/rs/zerolog"
	"os"
	"testing"
)

func TestLogging(t *testing.T) {
	Logger()
}

func TestDebuggingLogging(t *testing.T) {
	err := os.Setenv("LOG_LEVEL", "DEBUG")
	if err != nil {
		t.Failed()
	}
	logger := Logger()
	if logger.GetLevel() != zerolog.DebugLevel {
		t.Failed()
	}
}
