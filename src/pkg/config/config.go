package config

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"hms/gateway/pkg/common/utils"
	"os"
)

type Config struct {
	BaseUrl        string `json:"baseUrl"`
	DataPath       string `json:"dataPath"`
	Host           string `json:"host"`
	KeystoreKeyStr string `json:"keystoreKey"`
	KeystoreKey    []byte

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

	err = c.prepareKeystoreKey()

	return
}

// Make []byte key in KeystoreKey field from string in source configuration file
func (c *Config) prepareKeystoreKey() (err error) {
	keyByte, err := hex.DecodeString(c.KeystoreKeyStr)
	if err != nil {
		return
	}

	c.KeystoreKey = keyByte

	return
}
