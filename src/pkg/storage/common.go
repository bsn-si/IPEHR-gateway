package storage

import (
	"log"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/localfile"
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
