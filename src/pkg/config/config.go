package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common/utils"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type Config struct {
	BaseURL              string `json:"baseUrl"`
	DataPath             string `json:"dataPath"`
	Host                 string `json:"host"`
	KeystoreKey          string `json:"keystoreKey"`
	CreatingSystemID     string `json:"creatingSystemId"`
	CompressionEnabled   bool   `json:"compressionEnabled"`
	CompressionLevel     int    `json:"compressionLevel"` // 1-9 Fast-Best compression or 0 - No compression
	DefaultUserID        string `json:"defaultUserId"`
	DefaultGroupAccessID string `json:"defaultGroupAccessId"`
	Storage              struct {
		Localfile struct {
			Path string
		}
		Ipfs struct {
			EndpointURLs []string `json:"endpointURLs"`
		}
		Filecoin struct {
			LotusRPCEndpoint string
			BaseURL          string
			AuthToken        string
			DealsMaxPrice    uint64
			Miners           []string
		}
	}
	Contract struct {
		AddressEhrIndex    string
		AddressAccessStore string
		AddressUsers       string
		AddressDataStore   string
		Endpoint           string
		PrivKeyPath        string
		GasTipCap          int64 // maxPriorityFeePerGas used for hardhat testing
	}
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

	possibleConfigFiles := [4]string{
		userConfigFile,
		projectRootDir + "/" + mainConfigFile,
		os.Getenv("IPEHR_CONFIG_PATH"),
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
