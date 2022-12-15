package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Host    string
	BaseURL string
	LocalDB struct {
		Path       string
		Migrations string
	}
	Sync struct {
		Endpoint   string
		StartBlock uint64
		Contracts  []struct {
			Name    string
			Address string
			AbiPath string
		}
	}
}

const DefaultConfigPath = "config.json"

func New(path string) *Config {
	paths := []string{
		path,
		os.Getenv("IPEHR_STAT_CONFIG_PATH"),
		DefaultConfigPath,
	}

	cfg := Config{}

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

func (c *Config) load(path string) error {
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
