// TODO:
// validate config

package config

import (
	"log/slog" // avoiding logger package to prevent cyclic imports
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	e "github.com/tinkershack/meteomunch/errors"
)

type config struct {
	Munch          Munch     // Parameters of munch app, excluding external dependencies
	Mongo          DataStore // Gets picked if DocumentStore is "mongo"
	DLMRedis       DataStore
	MeteoProviders []MeteoProvider
}

func (c *config) GetMunch() Munch {
	return c.Munch
}

func (c *config) GetMongo() DataStore {
	return c.Mongo
}

func (c *config) GetDLMRedis() DataStore {
	return c.DLMRedis
}

func (c *config) GetMeteoProviders() []MeteoProvider {
	return c.MeteoProviders
}

// TODO: Validate URL string
type MeteoProvider struct {
	Name    string
	APIKey  string
	APIPath string // Path to the provider's API, excluding the base URI
	BaseURI string // URI of the provider's API, fully qualified with protocol
}

type Munch struct {
	Server   MunchServer
	LogLevel string // Log level for the application
}

type MunchServer struct {
	Hostname string
	Port     string
}

type DataStore struct {
	Name     string
	URI      string
	DBName   string
	DBNumber int
}

var defaultConfig = &config{
	Munch: Munch{
		Server: MunchServer{
			Hostname: "localhost",
			Port:     "50050",
		},
		LogLevel: "info",
	},
	Mongo: DataStore{
		Name:     "mongo",
		URI:      "mongodb://localhost:27017",
		DBName:   "meteomunch",
		DBNumber: 0,
	},
	DLMRedis: DataStore{
		Name:     "redis",
		URI:      "redis://localhost:6379",
		DBNumber: 1,
	},
	MeteoProviders: []MeteoProvider{
		{
			Name:    "open-meteo",
			APIKey:  "",
			APIPath: "v1/forecast",
			BaseURI: "https://api.open-meteo.com/",
		},
	},
}


// New creates a new config object
func New() (*config, error) {
	conf := defaultConfig

	pwd, err := os.Getwd()
	cobra.CheckErr(err)
	viper.AddConfigPath(pwd)
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.SetConfigName("munch")

	if err := viper.ReadInConfig(); err != nil {
		slog.Warn("Config file not found, using default config", "error", err)
		return defaultConfig, nil 
	}

	if err := viper.Unmarshal(conf); err != nil {
		slog.Error(e.FAIL, "err", err, "description", "Couldn't parse config, using default config")
		return defaultConfig, nil 
	}
	// Check for critical values and panic if missing
	validateCriticalFields(conf)

	return conf, nil
}
// validateCriticalFields checks for critical parameters in the configuration
func validateCriticalFields(conf *config) {
	for _, provider := range conf.MeteoProviders {
		if provider.Name == "meteoblue" && provider.APIKey == "" {
			slog.Error(e.FAIL, "error", "Critical config value missing: APIKey for MeteoProvider", "provider", provider.Name)
			panic("Critical config value is missing: API key for the provider")
		}
		if(provider.Name=="meteoblue" && provider.APIPath==""){
			slog.Error(e.FAIL,"error","Critical config value missing : API path for the meteoprovider","provider",provider.Name)
			panic("Critical config value is missing: APi path for the provider")
		}
		if(provider.Name=="meteoblue" && provider.BaseURI==""){
			slog.Error(e.FAIL,"error","Critical config vlaue is missing: BaseURI for the meteoprovider","provider",provider.Name)
			panic("Critical config value is missing : BaseURI for the provider")
		}
	}
}
