package storage

import (
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

		storage, err = localfile.Init(&cfg)
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
