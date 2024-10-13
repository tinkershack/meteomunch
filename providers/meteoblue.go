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
	client rest.HTTPClient
	config config.MeteoProvider
}

func NewMeteoBlueProvider() (*MeteoBlueProvider, error) {
	cfg, err := config.New()
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

func (p *MeteoBlueProvider) FetchData(qp map[string]string) (*plumber.BaseData, error) {
	// Creating a new request everytime we fecth data
	p.client.NewRequest()
	resp, err := p.client.SetQueryParams(qp).Get(p.config.APIPath)
	if err != nil {
		return nil, err
	}

	cfg, err := config.New()
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

	// Map the JSON response to plumber.BaseData
	baseData := &plumber.BaseData{
		Latitude:  data.Latitude,
		Longitude: data.Longitude,
		// Map other fields
	}

	return baseData, nil
}

// Maybe add a MeteoBlueQueryParams like OpenMeteoQueryParams to have some standard for all providers?
