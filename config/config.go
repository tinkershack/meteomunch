// TODO:
// validate config

package config

import (
	"log/slog" // avoiding logger package to prevent cyclic imports

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

// TODO: Validate URL sting
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

var legal = &struct {
	documentStore []string
}{
	documentStore: []string{"mongo"},
}

// TODO: This can be more perfomant by loading it just once and pass a copy for each call
func New() (*config, error) {
	// log := logger.NewTag("config")
	conf := new(config)

	if err := viper.Unmarshal(conf); err != nil {
		slog.Error(e.FAIL, "err", err, "description", "Couldn't parse config")
		return nil, err
	}

	return conf, nil
}
