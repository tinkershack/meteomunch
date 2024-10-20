package providers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tinkershack/meteomunch/config"
	"github.com/tinkershack/meteomunch/http/rest"
	"github.com/tinkershack/meteomunch/logger"
	"github.com/tinkershack/meteomunch/plumber"
)

const meteoBlueProviderName = "meteoblue"

type MeteoBlueProvider struct {
	client      rest.HTTPClient
	config      config.MeteoProvider
	queryParams map[string]string
	logLevel    string
}

func NewMeteoBlueProvider(cfg *config.Config) (*MeteoBlueProvider, error) {
	if cfg == nil {
		return nil, errors.New("configuration cannot be nil")
	}

	var meteoConfig config.MeteoProvider
	logLevel := cfg.Munch.LogLevel
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

	provider := MeteoBlueProvider{
		client:   client,
		config:   meteoConfig,
		logLevel: logLevel,
	}
	// Setting the default location to 0,0
	provider.SetQueryParams(plumber.NewCoordinates(0, 0))

	// Creating the new request(which will be reused in all FetchData calls) and setting the default queryParams on the client
	provider.client.NewRequest()
	provider.client.SetQueryParams(provider.queryParams)
	return &provider, nil
}

func (p *MeteoBlueProvider) FetchData(coords *plumber.Coordinates) (*plumber.BaseData, error) {
	resp, err := p.client.
		SetQueryParams(map[string]string{
			"lat": fmt.Sprintf("%f", coords.Latitude),
			"lon": fmt.Sprintf("%f", coords.Longitude),
		}).
		Get(p.config.APIPath)
	if err != nil {
		return nil, err
	}

	logger := logger.NewTag("providers:open-meteo")

	if p.logLevel == "debug" {
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
func (p *MeteoBlueProvider) SetQueryParams(coords *plumber.Coordinates) {
	// Setting the queryParams on the provider
	p.queryParams = map[string]string{
		"tz":            "GMT",
		"format":        "json",
		"forecast_days": "1",
		"apikey":        p.config.APIKey,
	}

}
