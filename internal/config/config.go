package config

import (
	"log"
)

type Configurator struct {
	config    Config
	mapConfig MapConfig
}

// NewConfigurator create new configurator
func NewConfigurator() Configurator {
	config, mapConfig := read()

	log.Println("successful application configuration")

	return Configurator{
		config:    config,
		mapConfig: mapConfig,
	}
}

// Get allows you to get structured config data
func (c Configurator) Get() Config {
	return c.config
}

// GetMap allows you to get non-structured config data
func (c Configurator) GetMap() MapConfig {
	return c.mapConfig
}
