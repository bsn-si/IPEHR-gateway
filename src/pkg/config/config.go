package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BaseUrl  string `json:"baseUrl"`
	DataPath string `json:"dataPath"`
	Host     string `json:"host"`

	path string
}

func New(configFilePath string) *Config {
	return &Config{
		path: configFilePath,
	}
}

func (c *Config) Reload() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}
