package config

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// Ensure the logger is set up correctly
	assert.NotNil(t, slog.Default())
}

func TestNewDefaultConfig(t *testing.T) {
	defaultConfig := NewDefaultConfig()
	assert.NotNil(t, defaultConfig)
	assert.Equal(t, "localhost", defaultConfig.Munch.Server.Hostname)
	assert.Equal(t, "50050", defaultConfig.Munch.Server.Port)
	assert.Equal(t, "info", defaultConfig.Munch.LogLevel)
	assert.Equal(t, "mongo", defaultConfig.Mongo.Name)
	assert.Equal(t, "mongodb://localhost:27017", defaultConfig.Mongo.URI)
	assert.Equal(t, "meteomunch", defaultConfig.Mongo.DBName)
	assert.Equal(t, 0, defaultConfig.Mongo.DBNumber)
	assert.Equal(t, "redis", defaultConfig.DLMRedis.Name)
	assert.Equal(t, "redis://localhost:6379", defaultConfig.DLMRedis.URI)
	assert.Equal(t, 1, defaultConfig.DLMRedis.DBNumber)
	assert.Equal(t, "open-meteo", defaultConfig.MeteoProviders[0].Name)
	assert.Equal(t, "", defaultConfig.MeteoProviders[0].APIKey)
	assert.Equal(t, "v1/forecast", defaultConfig.MeteoProviders[0].APIPath)
	assert.Equal(t, "https://api.open-meteo.com/", defaultConfig.MeteoProviders[0].BaseURI)
}

func TestGet(t *testing.T) {
	// Reset currentConfig for testing
	currentConfig = nil
	currentConfigErrors = nil

	config, err := Get()
	assert.NotNil(t, config)
	assert.Nil(t, err)
}

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "munch*.yml")
	assert.Nil(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`
Munch:
  LogLevel: "debug"
  Server:
    HostName: "localhost"
    Port: "12345"
MeteoProviders:
  - Name: "test-meteo"
    APIKey: "testkey"
    APIPath: "v1/test"
    BaseURI: "https://api.test-meteo.com"
Mongo:
  Name: "testmongo"
  URI: "mongodb://localhost:27017"
  DBName: "testdb"
  DBNumber: 2
DLMRedis:
  Name: "testredis"
  URI: "localhost:6379"
  DBNumber: 3
`)
	assert.Nil(t, err)
	tmpFile.Close()

	config, err := Load(nil, tmpFile.Name())
	assert.NotNil(t, config)
	assert.Nil(t, err)
	assert.Equal(t, "localhost", config.Munch.Server.Hostname)
	assert.Equal(t, "12345", config.Munch.Server.Port)
	assert.Equal(t, "debug", config.Munch.LogLevel)
	assert.Equal(t, "testmongo", config.Mongo.Name)
	assert.Equal(t, "mongodb://localhost:27017", config.Mongo.URI)
	assert.Equal(t, "testdb", config.Mongo.DBName)
	assert.Equal(t, 2, config.Mongo.DBNumber)
	assert.Equal(t, "testredis", config.DLMRedis.Name)
	assert.Equal(t, "localhost:6379", config.DLMRedis.URI)
	assert.Equal(t, 3, config.DLMRedis.DBNumber)
	assert.Equal(t, "test-meteo", config.MeteoProviders[0].Name)
	assert.Equal(t, "testkey", config.MeteoProviders[0].APIKey)
	assert.Equal(t, "v1/test", config.MeteoProviders[0].APIPath)
	assert.Equal(t, "https://api.test-meteo.com", config.MeteoProviders[0].BaseURI)
}

func TestValidateCriticalFields(t *testing.T) {
	config := NewDefaultConfig()
	config.MeteoProviders = append(config.MeteoProviders, MeteoProvider{
		Name:    "meteoblue",
		APIKey:  "",
		APIPath: "",
		BaseURI: "",
	})

	err := validateCriticalFields(config, false)
	assert.NotNil(t, err)
	assert.IsType(t, &CriticalErrors{}, err)
	criticalErrors := err.(*CriticalErrors)
	assert.Equal(t, 3, len(criticalErrors.Errors))
}
