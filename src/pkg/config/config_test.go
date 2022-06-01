package config

import "testing"

func Test_Config(t *testing.T) {

	t.Run("GetConfig fallback to example config file", func(t *testing.T) {
		defaultUrl := "http://localhost:8080"
		config := GetConfig("not-existed-file.json")
		if config.BaseUrl != defaultUrl {
			t.Errorf("BaseUrl mismatch. Expected %s, received %s", defaultUrl, config.BaseUrl)
		}
	})

}
