package config

import (
	"encoding/json"
	"errors"
)

// Config represents config for exporter.
type Config struct {
	Endpoint              string            `json:"endpoint"`
	Token                 string            `json:"token"`
	Queries               map[string]string `json:"queries"`
	RefreshDelaySeconds   int               `json:"refresh_delay_seconds"`
	RequestTimeoutSeconds int               `json:"request_timeout_seconds"`
	ListenPort            int               `json:"listen_port"`
}

const (
	defaultRequestTimeoutSeconds = 10
	defaultRefreshDelaySeconds   = 10
	defaultListenPort            = 8080
)

// New creates Config instance.
func New(raw []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(raw, &config)
	if err != nil {
		return nil, err
	}

	if config.Endpoint == "" {
		return nil, errors.New("empty endpoint")
	}

	if config.Token == "" {
		return nil, errors.New("empty token")
	}

	if len(config.Queries) == 0 {
		return nil, errors.New("empty queries")
	}

	if config.RequestTimeoutSeconds <= 0 {
		config.RequestTimeoutSeconds = defaultRequestTimeoutSeconds
	}

	if config.RefreshDelaySeconds <= 0 {
		config.RefreshDelaySeconds = defaultRefreshDelaySeconds
	}

	if config.ListenPort <= 0 {
		config.ListenPort = defaultListenPort
	}

	return &config, nil
}
