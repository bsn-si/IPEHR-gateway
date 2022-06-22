package config

import (
	"encoding/json"
	"errors"
	"hms/gateway/pkg/common/utils"
	"os"
)

type Config struct {
	BaseUrl            string `json:"baseUrl"`
	DataPath           string `json:"dataPath"`
	Host               string `json:"host"`
	KeystoreKey        string `json:"keystoreKey"`
	CompressionEnabled bool   `json:"compressionEnabled"`
	// 1-9 Fast-Best compression or 0 - No compression
	CompressionLevel int `json:"compressionLevel"`

	path string
}

var mainConfigFile = "config.json"
var fallbackConfigFile = "config.json.example"

func New(params ...string) (cfg *Config, err error) {
	var configFilePath string
	if len(params) == 1 {
		configFilePath = params[0]
	}
	path, err := resolveConfigFile(configFilePath)
	if err != nil {
		return
	}

	cfg = &Config{
		path: path,
	}
	err = cfg.load()

	return
}

// Resolves which config file we can use
func resolveConfigFile(userConfigFile string) (configFile string, err error) {
	projectRootDir, err := utils.ProjectRootDir()
	if err != nil {
		return
	}

	possibleConfigFiles := [3]string{
		userConfigFile,
		projectRootDir + "/" + mainConfigFile,
		projectRootDir + "/" + fallbackConfigFile,
	}

	for _, configFile = range possibleConfigFiles {
		_, err = os.Stat(configFile)
		if err == nil {
			return
		}
	}

	err = errors.New("not found any configuration file")

	return
}

// Loads content of source configuration file into configuration structure
func (c *Config) load() (err error) {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return
	}

	return
}
