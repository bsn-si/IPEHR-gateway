package config

import "testing"

func Test_Config(t *testing.T) {

	t.Run("GetConfig fallback to example config file", func(t *testing.T) {
		cfg := &Config{}
		err := Reload("../../../config.json.example", cfg)
		if err != nil {
			//
		}

		defaultUrl := "http://localhost:8080"
		if cfg.BaseUrl != defaultUrl {
			t.Errorf("BaseUrl mismatch. Expected %s, received %s", defaultUrl, cfg.BaseUrl)
		}
	})

}
