package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"hms/gateway/pkg/common/utils"
	"hms/gateway/pkg/errors"
)

type Config struct {
	BaseURL              string `json:"baseUrl"`
	DataPath             string `json:"dataPath"`
	Host                 string `json:"host"`
	StoragePath          string `json:"storagePath"`
	KeystoreKey          string `json:"keystoreKey"`
	CreatingSystemID     string `json:"creatingSystemId"`
	CompressionEnabled   bool   `json:"compressionEnabled"`
	CompressionLevel     int    `json:"compressionLevel"` // 1-9 Fast-Best compression or 0 - No compression
	DefaultUserID        string `json:"defaultUserId"`
	DefaultGroupAccessID string `json:"defaultGroupAccessId"`
	Storage              struct {
		Ipfs struct {
			EndpointURL string `json:"endpointUrl"`
		} `json:"ipfs"`
	} `json:"storage"`
	Contract struct {
		Address     string `json:"address"`
		Endpoint    string `json:"endpoint"`
		PrivKeyPath string `json:"privKeyPath"`
	} `json:"contract"`
	DB struct {
		FilePath string `json:"filePath"`
	} `json:"db"`

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

	cfgJSON, _ := json.MarshalIndent(cfg, "", "    ")

	log.Println("IPEHR Config:", string(cfgJSON))

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
		if _, err = os.Stat(configFile); err == nil {
			return configFile, nil
		}
	}

	return "", fmt.Errorf("Not found any configuration file: %w", errors.ErrIsNotExist)
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
