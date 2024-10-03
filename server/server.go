package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/tinkershack/meteomunch/config"
	e "github.com/tinkershack/meteomunch/errors"
	"github.com/tinkershack/meteomunch/http/rest"
	"github.com/tinkershack/meteomunch/logger"
)

func Serve(ctx context.Context, args []string) {
	logger := logger.NewTag("server")

	// The following are transient test routes for an early stage development convenience.
	// These will be cleanedup in the future.

	cfg, err := config.New()
	if err != nil {
		logger.Error(e.FAIL, "err", err, "description", "Couldn't parse config")
		os.Exit(-1)
	}
	logger.Info("Config parsed successfully", "config", cfg)

	var client rest.HTTPClient = rest.NewClient().
		SetAuthToken("dummy-auth-token").
		// FIX: SetQueryParams() isn't getting picked up
		SetQueryParams(map[string]string{
			"lat":           "47.558",
			"lon":           "7.587",
			"asl":           "279",
			"tz":            "Europe/Zurich",
			"name":          "Test",
			"windspeed":     "kmh",
			"format":        "json",
			"history_days":  "1",
			"forecast_days": "0",
			"apikey":        cfg.MeteoProviders[0].APIKey,
		}).
		AcceptJSON().
		EnableTrace().
		SetDefaults().
		SetBaseURL(cfg.MeteoProviders[0].URI).
		SetDebug()
		// FIX: SetQueryString isn't getting picked up
		// SetQueryString("lat=47.558&lon=7.573&asl=279&tz=Europe%2FZurich&name=Test&windspeed=kmh&format=json&history_days=1&forecast_days=0&apikey=" + cfg.MeteoProviders[0].APIKey)

	qs := "lat=47.558&lon=7.573&asl=279&tz=Europe%2FZurich&name=Test&windspeed=kmh&format=json&history_days=1&forecast_days=0&apikey=" + cfg.MeteoProviders[0].APIKey
	uri := fmt.Sprintf("packages/air-1h_air-day?%s", qs)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
		logger.Debug(r.URL.String())
	})
	mux.HandleFunc("GET /meteo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
		logger.Info(r.URL.String())

		resp, err := client.Get(uri)
		if err != nil {
			logger.Error(e.FAIL, "err", err, "description", "HTTP request failed!")
		}

		fmt.Println("Response Status:", resp.Status())
		fmt.Println("Response Body:", string(resp.Body()))

		traceInfo := resp.TraceInfo()
		fmt.Printf("Trace Info: %+v\n", traceInfo)
	})

	logger.Info("Ready, Plank? Serving Meteo Munch on " + cfg.Munch.Server.Hostname + ":" + cfg.Munch.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Munch.Server.Hostname, cfg.Munch.Server.Port), mux)
	logger.Error(e.FATAL, "err", err, "description", "Server killed!")
	os.Exit(-1)
}
