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

const openMeteoProviderName = "open-meteo"

type OpenMeteoProvider struct {
	client rest.HTTPClient
	config config.MeteoProvider
}

// NewOpenMeteoProvider returns a new instance of OpenMeteoProvider
func NewOpenMeteoProvider() (*OpenMeteoProvider, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	var meteoConfig config.MeteoProvider
	found := false

	for _, provider := range cfg.MeteoProviders {
		if provider.Name == openMeteoProviderName {
			meteoConfig = provider
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("open-meteo provider configuration not found")
	}

	client := rest.NewClient().SetDefaults().SetBaseURL(meteoConfig.BaseURI)
	if cfg.Munch.LogLevel == "debug" {
		client.SetDebug()
		client.EnableTrace()
	}

	return &OpenMeteoProvider{
		client: client,
		config: meteoConfig,
	}, nil
}

// FetchData fetches API data from open-meteo provider for the given query parameters map
func (p *OpenMeteoProvider) FetchData(coords *plumber.Coordinates) (*plumber.BaseData, error) {
	p.client.NewRequest()
	qp := p.SetQueryParams(coords)

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

	data := new(plumber.BaseData) // Fields of the struct will be zero-initialized

	if err := json.Unmarshal(resp.Body(), data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
}

// SetQueryParams forms the query parameters for OpenMeteo API based on given coordinates
func (p *OpenMeteoProvider) SetQueryParams(coords *plumber.Coordinates) map[string]string {
	return map[string]string{
		"latitude":       fmt.Sprintf("%f", coords.Latitude),
		"longitude":      fmt.Sprintf("%f", coords.Longitude),
		"current":        "temperature_2m,relative_humidity_2m,apparent_temperature,is_day,precipitation,rain,showers,snowfall,weather_code,cloud_cover,pressure_msl,surface_pressure,wind_speed_10m,wind_direction_10m,wind_gusts_10m",
		"hourly":         "temperature_2m,relative_humidity_2m,dew_point_2m,apparent_temperature,precipitation_probability,precipitation,weather_code,pressure_msl,surface_pressure,cloud_cover,cloud_cover_low,cloud_cover_mid,cloud_cover_high,visibility,evapotranspiration,et0_fao_evapotranspiration,vapour_pressure_deficit,wind_speed_10m,wind_speed_80m,wind_speed_120m,wind_speed_180m,wind_direction_10m,wind_direction_80m,wind_direction_120m,wind_direction_180m,wind_gusts_10m,temperature_80m,temperature_120m,temperature_180m,uv_index,uv_index_clear_sky,is_day,sunshine_duration,total_column_integrated_water_vapour,cape,lifted_index,convective_inhibition,freezing_level_height,boundary_layer_height,temperature_1000hPa,temperature_975hPa,temperature_950hPa,temperature_925hPa,temperature_900hPa,temperature_850hPa,temperature_800hPa,temperature_700hPa,temperature_600hPa,temperature_500hPa,temperature_400hPa,relative_humidity_1000hPa,relative_humidity_975hPa,relative_humidity_950hPa,relative_humidity_925hPa,relative_humidity_900hPa,relative_humidity_850hPa,relative_humidity_800hPa,relative_humidity_700hPa,relative_humidity_600hPa,relative_humidity_500hPa,relative_humidity_400hPa,cloud_cover_1000hPa,cloud_cover_975hPa,cloud_cover_950hPa,cloud_cover_925hPa,cloud_cover_900hPa,cloud_cover_850hPa,cloud_cover_800hPa,cloud_cover_700hPa,cloud_cover_600hPa,cloud_cover_500hPa,cloud_cover_400hPa,wind_speed_1000hPa,wind_speed_975hPa,wind_speed_950hPa,wind_speed_925hPa,wind_speed_900hPa,wind_speed_850hPa,wind_speed_800hPa,wind_speed_700hPa,wind_speed_600hPa,wind_speed_500hPa,wind_speed_400hPa,wind_direction_1000hPa,wind_direction_975hPa,wind_direction_950hPa,wind_direction_925hPa,wind_direction_900hPa,wind_direction_850hPa,wind_direction_800hPa,wind_direction_700hPa,wind_direction_600hPa,wind_direction_500hPa,wind_direction_400hPa,geopotential_height_1000hPa,geopotential_height_975hPa,geopotential_height_950hPa,geopotential_height_925hPa,geopotential_height_900hPa,geopotential_height_850hPa,geopotential_height_800hPa,geopotential_height_700hPa,geopotential_height_600hPa,geopotential_height_500hPa,geopotential_height_400hPa",
		"daily":          "weather_code,temperature_2m_max,temperature_2m_min,apparent_temperature_max,apparent_temperature_min,sunrise,sunset,daylight_duration,sunshine_duration,uv_index_max,uv_index_clear_sky_max,precipitation_sum,precipitation_hours,precipitation_probability_max,wind_speed_10m_max,wind_gusts_10m_max,wind_direction_10m_dominant,shortwave_radiation_sum,et0_fao_evapotranspiration",
		"timeformat":     "unixtime",
		"timezone":       "GMT",
		"forecast_days":  "1",
		"forecast_hours": "24",
		"cell_selection": "nearest",
		"models":         "best_match",
	}
}
