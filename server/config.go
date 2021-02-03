package server

import (
	"encoding/json"
	"io/ioutil"
)

// Config stores the server configuration.
type Config struct {
	LogPath    string `json:"logPath"`
	SecretPath string `json:"secretPath"`
}

// ReadConfig reads the JSON config from `configPath`.
func ReadConfig(configPath string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err := json.Unmarshal(configBytes, config); err != nil {
		return nil, err
	}

	if config.LogPath == "" {
		config.LogPath = "/var/log/simple-server.log"
	}

	if config.SecretPath == "" {
		config.SecretPath = "secret.key"
	}

	return config, nil
}
