// Package providers offers an interface and a factory method for weather data providers.
// This package is designed to facilitate the integration of various weather data providers
// by defining a common interface that each provider must implement. It also includes a
// factory method to instantiate the appropriate provider based on a given name.
//
// The package relies on the following external packages:
// - github.com/tinkershack/meteomunch/logger: For logging purposes.
// - github.com/tinkershack/meteomunch/plumber: For handling the base data structure.
//
// The main components of this package are:
// - Provider interface: Defines the methods that each provider must implement.
// - NewProvider function: A factory method that returns the appropriate provider based on the name.
//
// Example usage:
//
//	provider, err := providers.NewProvider("open-meteo")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	data, err := provider.FetchData(queryParams)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// The package initializes a logger with the tag "providers" to facilitate logging within the package.
package providers

import (
	"fmt"

	"github.com/tinkershack/meteomunch/config"
	"github.com/tinkershack/meteomunch/logger"
	"github.com/tinkershack/meteomunch/plumber"
)

var l logger.Logger

func init() {
	l = logger.NewTag("providers")
}

// Provider interface defines the methods that each provider must implement
type Provider interface {
	FetchData(coords *plumber.Coordinates) (*plumber.BaseData, error)
	SetQueryParams(coords *plumber.Coordinates)
}

// NewProvider returns the appropriate provider based on the name
func NewProvider(name string, cfg *config.Config) (Provider, error) {
	switch name {
	case "open-meteo":
		p, err := NewOpenMeteoProvider(cfg)
		if err != nil {
			return nil, err
		}
		return p, nil
	case "meteoblue":
		p, err := NewMeteoBlueProvider(cfg)
		if err != nil {
			return nil, err
		}
		return p, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
