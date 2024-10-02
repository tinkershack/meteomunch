// TODO:
// validate config

package config

import (
	"github.com/spf13/viper"
	e "github.com/tinkershack/meteomunch/errors"
	"github.com/tinkershack/meteomunch/logger"
)

type config struct {
	Munch         munch     // Paramters of munch app, excluding external dependencies
	DocumentStore string    // Name of the prefered DS
	Mongo         dataStore // Gets picked if DocumentStore is "mongo"
	DLMRedis      dataStore
}

type munch struct {
	Serve munchServer
}

type munchServer struct {
	Hostname string
	Port     string
}

type dataStore struct {
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

func New() (*config, error) {
	log := logger.NewTag("config")
	conf := new(config)

	if err := viper.Unmarshal(conf); err != nil {
		log.Error(e.FAIL, "err", err, "description", "Couldn't parse config")
		return nil, err
	}

	return conf, nil
}
