package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/tinkershack/meteomunch/config"
	e "github.com/tinkershack/meteomunch/errors"
	"github.com/tinkershack/meteomunch/http/rest"
	"github.com/tinkershack/meteomunch/logger"
	"github.com/tinkershack/meteomunch/plumber"
)

const meteoBlueProviderName = "meteoblue"

type MeteoBlueProvider struct {
	client rest.HTTPClient
	config config.MeteoProvider
}

func NewMeteoBlueProvider() (*MeteoBlueProvider, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	var meteoConfig config.MeteoProvider
	found := false

	for _, provider := range cfg.MeteoProviders {
		if provider.Name == meteoBlueProviderName {
			meteoConfig = provider
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("meteoblue provider configuration not found")
	}

	client := rest.NewClient().SetDefaults().SetBaseURL(meteoConfig.BaseURI)
	if cfg.Munch.LogLevel == "debug" {
		client.SetDebug()
		client.EnableTrace()
	}

	return &MeteoBlueProvider{
		client: client,
		config: meteoConfig,
	}, nil
}

func (p *MeteoBlueProvider) FetchData(coords *plumber.Coordinates) (*plumber.BaseData, error) {
	// Creating a new request everytime we fecth data
	p.client.NewRequest()
	// Getting the query parameters
	qp := p.SetQueryParams(coords)

	// Setting the paramters on the request and making the request
	resp, err := p.client.SetQueryParams(qp).Get(p.config.APIPath)
	if err != nil {
		return nil, err
	}

	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	logger := logger.NewTag("providers:open-meteo")

	if cfg.Munch.LogLevel == "debug" {
		// Refraining from logger.Debug() as it doesn't pretty print the resty response stats and body
		// And, this is just for debugging purposes so it's okay to use fmt.Println as it's not logged in production
		logger.Debug("Response", "status:", resp.Status())
		logger.Debug("Response", "body:", string(resp.Body()))

		traceInfo := resp.TraceInfo()
		logger.Debug("Response", "trace", fmt.Sprintf("%+v", traceInfo))
	}

	var data struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
		// Add other fields as per the API response
	}

	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// TO-DO: #12
	// Map the JSON response to plumber.BaseData
	baseData := &plumber.BaseData{
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
		// Map other fields
	}

	return baseData, nil
}

// Set the QueryParams for the request
func (p *MeteoBlueProvider) SetQueryParams(coords *plumber.Coordinates) map[string]string {
	// Get the config to accesss the API Key
	cfg, err := config.Get()
	if err != nil {
		slog.Error("Couldn't parse config", "error", err, "description", e.FAIL)
		os.Exit(-1)
	}
	return map[string]string{
		"lat":           fmt.Sprintf("%f", coords.Latitude),
		"lon":           fmt.Sprintf("%f", coords.Longitude),
		"tz":            "GMT",
		"format":        "json",
		"forecast_days": "1",
		"apikey":        cfg.MeteoProviders[0].APIKey,
	}

}
