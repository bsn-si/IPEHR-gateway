package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/service/syncer"
)

type StatConfig struct {
	Host    string
	BaseURL string
	LocalDB struct {
		Path       string
		Migrations string
	}
	Sync syncer.Config
}

const DefaultConfigPath = "config.json"

func NewStatConfig(path string) *StatConfig {
	paths := []string{
		path,
		os.Getenv("IPEHR_STAT_CONFIG_PATH"),
		DefaultConfigPath,
	}

	cfg := StatConfig{}

	for _, path = range paths {
		err := cfg.load(path)
		if err == nil {
			cfgJSON, _ := json.MarshalIndent(cfg, "", "    ")
			log.Println("ipehr-stat Config:", string(cfgJSON))
			return &cfg
		}
	}

	log.Fatal("No suitable config found")

	return nil
}

func (c *StatConfig) load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Load config '%s' error: %v", path, err)
		return err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		log.Printf("Unmarshal config '%s' error: %v", path, err)
		return err
	}

	return nil
}
