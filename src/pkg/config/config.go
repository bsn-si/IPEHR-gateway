package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BaseUrl     string `json:"baseUrl"`
	DataPath    string `json:"dataPath"`
	Host        string `json:"host"`
	StoragePath string `json:"storagePath"`

	path string
}

func New(configFilePath string) *Config {
	return &Config{
		path:        configFilePath,
		StoragePath: "",
	}
}

func (c *Config) Reload() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}
