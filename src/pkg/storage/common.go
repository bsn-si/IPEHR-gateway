package storage

import (
	config2 "hms/gateway/pkg/config"
	"hms/gateway/pkg/storage/localfile"
	"log"
)

var storage Storager

func Init(sc *Config) {
	if storage == nil {
		cfg := localfile.Config{
			BasePath: sc.Path(),
			Depth:    3,
		}

		var err error

		globalConfig, err := config2.New()
		if err != nil {
			log.Fatal(err)
		}

		storage, err = localfile.Init(&cfg, globalConfig)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Storage() Storager {
	if storage == nil {
		log.Fatal("Storage is not initialized")
	}

	return storage
}
