package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}
