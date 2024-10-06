package logger

import (
	"log/slog"
	"os"

	"github.com/tinkershack/meteomunch/config"
	e "github.com/tinkershack/meteomunch/errors"
)

type Logger *slog.Logger

// NewTag returns a new instance of a logger instance that sticks a string tag key-value pair to every log.
func NewTag(tag string) *slog.Logger {
	logger := New().With("tag", tag)
	return logger
}

// New returns a new generic instance of a logger that implements Logger interface
func New() *slog.Logger {
	cfg, err := config.New()
	if err != nil {
		slog.Error(e.FAIL, "err", err, "description", "Couldn't parse config")
	}

	var level slog.Level
	switch cfg.Munch.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	// slog.SetDefault(logger)
	return logger
}
