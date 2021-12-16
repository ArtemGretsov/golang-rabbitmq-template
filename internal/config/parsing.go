package config

import (
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/ArtemGretsov/golang-rabbitmq-template/config"
)

const (
	DefaultPath = "config/default.json"
	LocalPath   = "config/local.json"
)

type Config configtype.Config
type MapConfig map[string]interface{}

func read() (Config, MapConfig) {
	result := make(map[string]interface{})
	conf := Config{}

	if err := readJSONConfig(DefaultPath, result, true); err != nil {
		log.Fatalln(err)
	}

	readEnv(result)

	if err := readJSONConfig(LocalPath, result, false); err != nil {
		log.Fatalln(err)
	}

	if err := mapstructure.Decode(result, &conf); err != nil {
		log.Fatalln(err)
	}

	return conf, result
}
