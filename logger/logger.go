package logger

import (
	"log/slog"
	"os"
)

type Logger *slog.Logger

// NewTag returns a new instance of a logger instance that sticks a string tag key-value pair to every log.
func NewTag(tag string) *slog.Logger {
	logger := New().With("tag", tag)
	return logger
}

// New returns a new generic instance of a logger that implements Logger interface
func New() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// slog.SetDefault(logger)
	return logger
}
