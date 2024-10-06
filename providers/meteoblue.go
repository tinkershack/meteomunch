package providers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tinkershack/meteomunch/config"
	"github.com/tinkershack/meteomunch/http/rest"
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

	return &MeteoBlueProvider{
		client: rest.NewClient(),
		config: meteoConfig,
	}, nil
}

func (p *MeteoBlueProvider) FetchData(qp map[string]string) (*plumber.BaseData, error) {
	resp, err := p.client.SetQueryParams(map[string]string{
		"lat": "11.25",
		"lon": "77",
	}).Get("https://my.meteoblue.com/v1/forecast")
	if err != nil {
		return nil, err
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
