package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type LoadConfigOptions struct {
	ConfigFile string
}

func LoadConfig(options *LoadConfigOptions) (*BrokeConfig, error) {
	if !strings.HasSuffix(options.ConfigFile, ".broke.yml") {
		errorMessage := fmt.Sprintf("Invalid file extension on input config file '%s'. Must be *.broke.yml", options.ConfigFile)
		log.Error().Msg(errorMessage)
		return nil, errors.New(errorMessage)
	}

	file, err := os.Open(options.ConfigFile)
	if err != nil {
		errorMessage := "Failed to open configuration file: " + options.ConfigFile
		log.Error().Err(err).Msg(errorMessage)
		return nil, errors.New(errorMessage)
	}
	defer file.Close()

	var config BrokeConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		errorMessage := "Failed to parse configuration file: " + options.ConfigFile
		log.Error().Err(err).Msg(errorMessage)
		return nil, errors.New(errorMessage)
	}

	// TODO: make some validations:
	// - User Sources and User Targets must have unique names

	return &config, nil
}
