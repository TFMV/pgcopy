package model

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// DBConfig represents the configuration for a database connection.
type DBConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Pass         string `yaml:"pass"`
	DB           string `yaml:"db"`
	IsUnixSocket bool   `yaml:"isUnixSocket"`
}

// Config represents the overall configuration including source and target databases and tables to be copied.
type Config struct {
	Source DBConfig `yaml:"source"`
	Target DBConfig `yaml:"target"`
	Tables []string `yaml:"tables"`
}

// JsonResponse represents the structure of the JSON response.
type JsonResponse struct {
	Message   string  `json:"message"`
	TimeTaken float64 `json:"timeTaken"`
}

// GetConf reads the YAML configuration file and unmarshals it into a Config struct.
func GetConf(filePath string) (*Config, error) {
	var config Config
	yamlContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
