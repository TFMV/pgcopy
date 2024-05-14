package config

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Threads        int
	BatchSize      int
	Query          string
	SDBSystem      string
	SDBName        string
	SDBUser        string
	SDBPassword    string
	TDBSystem      string
	TTable         string
	TDBName        string
	TDBUser        string
	TDBPassword    string
	ConnectionType string
	UnixSocketPath string
}

func LoadConfig(cfg Config) (Config, error) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
