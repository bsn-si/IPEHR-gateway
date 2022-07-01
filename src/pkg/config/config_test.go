package config_test

import (
	"testing"

	"hms/gateway/pkg/config"
)

func Test_Config(t *testing.T) {
	t.Run("GetConfig fallback to example config file", func(t *testing.T) {
		cfg, err := config.New()
		if err != nil {
			t.Fatal(err)
		}

		defaultURL := "http://localhost:8080"
		if cfg.BaseURL != defaultURL {
			t.Errorf("BaseUrl mismatch. Expected %s, received %s", defaultURL, cfg.BaseURL)
		}
	})
}
