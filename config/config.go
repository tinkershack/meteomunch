// Package config provides functionality to load and validate the configuration
// for the meteomunch application. It supports loading configuration from
// default values or from a configuration file. The package also includes
// validation for critical configuration fields to ensure the application
// runs with the necessary parameters.
//
// Note: The functions config.Get and config.Load are not safe for concurrent use.
// Ensure that these functions are called in a single-threaded context or
// protected by appropriate synchronization mechanisms if used concurrently.
package config

import (
	"log/slog" // avoiding logger package to prevent cyclic imports
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	e "github.com/tinkershack/meteomunch/errors"
)

func init() {
	// Configure slog to use JSON format
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("tag", "config")
	slog.SetDefault(logger)
}

type Config struct {
	Munch          Munch     // Parameters of munch app, excluding external dependencies
	Mongo          DataStore // Gets picked if DocumentStore is "mongo"
	DLMRedis       DataStore
	MeteoProviders []MeteoProvider
}

func (c *Config) GetMunch() Munch {
	return c.Munch
}

func (c *Config) GetMongo() DataStore {
	return c.Mongo
}

func (c *Config) GetDLMRedis() DataStore {
	return c.DLMRedis
}

func (c *Config) GetMeteoProviders() []MeteoProvider {
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

var currentConfig *Config
var currentConfigErrors error

var defaultConfig = &Config{
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

// NewDefaultConfig returns a deep copy of the default configuration
//
// By doing this, we ensure that newConfig has its own independent copy of the MeteoProviders slice. Therefore, any changes made to newConfig will not affect defaultConfig.
func NewDefaultConfig() *Config {
	newConfig := *defaultConfig
	newConfig.MeteoProviders = make([]MeteoProvider, len(defaultConfig.MeteoProviders))
	copy(newConfig.MeteoProviders, defaultConfig.MeteoProviders)
	return &newConfig
}

// Get returns the current configuration that was last loaded, or loads the configuration if it hasn't been loaded yet
func Get() (*Config, error) {
	if currentConfig == nil {
		currentConfig, currentConfigErrors = Load(nil, "")
		// return cfg, err
	}
	return currentConfig, currentConfigErrors
}

// Load loads the configuration from the file if a preferred config is not provided
//
// This is useful when Munch is being used as a package and the user wants to dynamically pass their own configuration
func Load(c *Config, configPath string) (*Config, error) {
	currentConfig = NewDefaultConfig() // Assign it to config.currentConfig

	// Use the provided config if it's not nil
	if c != nil {
		currentConfig = c
	}

	// Load the config from the file if it hasn't been provided
	if c == nil {
		pwd, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(pwd)
		viper.AddConfigPath(".")
		viper.SetConfigType("yml")
		viper.SetConfigName("munch")
		if configPath != "" {
			viper.SetConfigFile(configPath)
		}

		slog.Info("Reading config file", "path", viper.ConfigFileUsed())
		if err := viper.ReadInConfig(); err != nil {
			slog.Error("Config file not found, using default config", "error", err)
			return currentConfig, nil
		}

		if err := viper.Unmarshal(currentConfig); err != nil {
			slog.Error(e.FAIL, "err", err, "description", "Couldn't parse config, using default config")
			return currentConfig, nil
		}
	}

	// Check for critical values, but do not panic. Instead, return an error
	if err := validateCriticalFields(currentConfig, false); err != nil {
		currentConfigErrors = err
		return currentConfig, currentConfigErrors
	}

	currentConfigErrors = nil // Reset the errors
	return currentConfig, currentConfigErrors
}

// CriticalError represents a custom error type for critical config errors
type CriticalError struct {
	Field   string
	Message string
}

// Error implements the error interface for CriticalError
func (e *CriticalError) Error() string {
	return "Critical config value missing: " + e.Field + " - " + e.Message
}

// CriticalErrors represents a collection of critical validation errors
type CriticalErrors struct {
	Errors []error
}

// Error implements the error interface for CriticalValidationErrors
func (e *CriticalErrors) Error() string {
	var errMsg string
	for _, err := range e.Errors {
		errMsg += err.Error()
	}
	return errMsg
}

// validateCriticalFields checks for critical parameters in the configuration
func validateCriticalFields(conf *Config, isPanic bool) error {
	var ve []error // Validation errors

	for _, provider := range conf.MeteoProviders {
		if provider.Name == "meteoblue" {
			if provider.APIKey == "" {
				ve = append(ve, &CriticalError{
					Field:   "APIKey",
					Message: provider.Name,
				})
			}
			if provider.APIPath == "" {
				ve = append(ve, &CriticalError{
					Field:   "APIPath",
					Message: provider.Name,
				})
			}
			if provider.BaseURI == "" {
				ve = append(ve, &CriticalError{
					Field:   "BaseURI",
					Message: provider.Name,
				})
			}
		}
	}

	if len(ve) > 0 {
		if isPanic {
			for _, err := range ve {
				slog.Error(e.FAIL, "error", err.Error())
			}
			panic("Critical config values are missing")
		}
		return &CriticalErrors{Errors: ve}
	}

	return nil
}
