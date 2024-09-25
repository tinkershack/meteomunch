package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/tinkershack/meteomunch/logger"
)

func main() {
	logger := logger.New()
	slog.SetDefault(logger)
	logger.Info("Ready, Plank? Serving Meteo Munch.")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	err := http.ListenAndServe("localhost:8080", mux)
	logger.Error("fatal", "err", err, slog.Int("error_code", -1), "description", "Server killed!")
	os.Exit(-1)
}
