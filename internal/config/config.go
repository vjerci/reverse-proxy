package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var ErrConfigNotSet = errors.New("CONFIG_FILE env var not set")
var ErrConfigRead = errors.New("couldn't read config file")
var ErrConfigJSON = errors.New("couldn't decode json of config file")

type ConfigData struct {
	ForwardHost   string          `json:"forward_host"`
	ForwardScheme string          `json:"forward_scheme"`
	Block         [][]interface{} `json:"block"`
}

func Load(configFilePath string) (*ConfigData, error) {
	if configFilePath == "" {
		return nil, ErrConfigNotSet
	}

	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigRead, err)
	}

	var configData ConfigData

	err = json.NewDecoder(bytes.NewBuffer(fileContent)).Decode(&configData)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigJSON, err)
	}

	return &configData, nil
}
