package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Config struct {
	BaseUrl string `json:"baseUrl"`
}

var (
	config Config
	once   sync.Once
)

const defaultConfigPath = "../../../config.json.example"

func GetConfig(configFileName string) *Config {
	once.Do(func() {
		lookupConfig(configFileName, &config)
	})

	return &config
}

func lookupConfig(configFileName string, configuration interface{}) *Config {
	wd, _ := os.Getwd()
	configFilePath := wd + "/" + configFileName

	err := readJson(configFilePath, &config)
	if err != nil && os.IsNotExist(err) {
		err = readJson(wd+"/"+defaultConfigPath, &config)
	}

	if err != nil {
		log.Println("Error reading config: ", err)
	}

	return &config
}

func readJson(configFile string, configuration interface{}) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &configuration)
}
