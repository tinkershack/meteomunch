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
	DocumentStore  string    // Name of the preferred DS
	Mongo          DataStore // Gets picked if DocumentStore is "mongo"
	DLMRedis       DataStore
	MeteoProviders []MeteoProvider
}

func (c *config) GetMunch() Munch {
	return c.Munch
}

func (c *config) GetDocumentStore() string {
	return c.DocumentStore
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
	DocumentStore: "mongo",
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

var legal = &struct {
	documentStore []string
}{
	documentStore: []string{"mongo"},
}

// New creates a new config object
func New() (*config, error) {
	conf := new(config)

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

	setDefaultValues(conf)

	// Check for critical values and panic if missing
	validateCriticalFields(conf)

	return conf, nil
}

// setDefaultValues sets default values for non-critical parameters in the configuration
func setDefaultValues(conf *config) {
	//default values if not present
	if conf.Munch.Server.Hostname == "" {
		conf.Munch.Server.Hostname = defaultConfig.Munch.Server.Hostname
		slog.Warn("Hostname missing, using default", "value", defaultConfig.Munch.Server.Hostname)
	}

	if conf.Munch.Server.Port == "" {
		conf.Munch.Server.Port = defaultConfig.Munch.Server.Port
		slog.Warn("Port missing, using default", "value", defaultConfig.Munch.Server.Port)
	}

	if conf.Munch.LogLevel == "" {
		conf.Munch.LogLevel = defaultConfig.Munch.LogLevel
		slog.Warn("LogLevel missing, using default", "value", defaultConfig.Munch.LogLevel)
	}

	if conf.DocumentStore == "" {
		conf.DocumentStore = defaultConfig.DocumentStore
		slog.Warn("DocumentStore missing, using default", "value", defaultConfig.DocumentStore)
	}

	//default values for Mongo configuration if missing
	if conf.Mongo.URI == "" {
		conf.Mongo.URI = defaultConfig.Mongo.URI
		slog.Warn("Mongo URI missing, using default", "value", defaultConfig.Mongo.URI)
	}

	if conf.Mongo.DBName == "" {
		conf.Mongo.DBName = defaultConfig.Mongo.DBName
		slog.Warn("Mongo DBName missing, using default", "value", defaultConfig.Mongo.DBName)
	}

	if(conf.MeteoProviders[0].Name==""){
		conf.MeteoProviders[0].Name=defaultConfig.MeteoProviders[0].Name
	}
	if(conf.MeteoProviders[0].APIPath==""){
		if(conf.MeteoProviders[0].Name=="open-meteo"){
			conf.MeteoProviders[0].APIPath=defaultConfig.MeteoProviders[0].APIPath
		}else { 
			validateCriticalFields(conf);
		}
	}
	if(conf.MeteoProviders[0].BaseURI==""){
		if(conf.MeteoProviders[0].Name=="open-meteo"){
			conf.MeteoProviders[0].BaseURI=defaultConfig.MeteoProviders[0].BaseURI
		}else{
			validateCriticalFields(conf);
		}
	}
	// Check for non-critical MeteoProviders and set defaults
	if len(conf.MeteoProviders) == 0 {
		slog.Warn("No MeteoProviders configured, using default config")
		conf.MeteoProviders = defaultConfig.MeteoProviders
	}
}

// validateCriticalFields checks for critical parameters in the configuration
func validateCriticalFields(conf *config) {
	for _, provider := range conf.MeteoProviders {
		if provider.Name != "open-meteo" && provider.APIKey == "" {
			slog.Error(e.FAIL, "error", "Critical config value missing: APIKey for MeteoProvider", "provider", provider.Name)
			panic("Critical config value is missing: API key for the provider")
		}
		if(provider.Name!="open-meteo" && provider.APIPath==""){
			slog.Error(e.FAIL,"error","Critical config value missing : API path for the meteoprovider","provider",provider.Name)
			panic("Critical config value is missing: APi path for the provider")
		}
		if(provider.Name!="open-meteo" && provider.BaseURI==""){
			slog.Error(e.FAIL,"error","Critical config vlaue is missing: BaseURI for the meteoprovider","provider",provider.Name)
			panic("Critical config value is missing : BaseURI for the provider")
		}
	}
}
