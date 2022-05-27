package storage

import (
	"hms/gateway/pkg/storage/localfile"
	"log"
	"os"
	"path/filepath"
)

var Storage Storager

func Init() Storager {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(ex)

	//TODO getting basepath from general config

	if Storage == nil {
		cfg := localfile.Config{
			BasePath: path,
			Depth:    3,
		}
		var err error
		Storage, err = localfile.Init(&cfg)
		if err != nil {
			log.Fatal(err)
			return nil
		}
	}
	return Storage
}
