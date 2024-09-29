package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	e "github.com/tinkershack/meteomunch/errors"
	"github.com/tinkershack/meteomunch/logger"
)

func Serve(ctx context.Context, args []string) {
	logger := logger.NewTag("server")
	logger.Info("Ready, Plank? Serving Meteo Munch.")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
		logger.Debug(r.URL.String())
	})

	err := http.ListenAndServe("localhost:8080", mux)
	logger.Error(e.FATAL, "err", err, "description", "Server killed!")
	os.Exit(-1)
}
