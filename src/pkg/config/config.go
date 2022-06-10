package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BaseUrl  string `json:"baseUrl"`
	DataPath string `json:"dataPath"`
	Host     string `json:"host"`
}

func Reload(configFilePath string, cfg *Config) error {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, cfg)
}
