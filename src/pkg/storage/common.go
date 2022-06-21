package storage

import (
	"hms/gateway/pkg/storage/localfile"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Storage Storager

func getStorageName() string {
	name := os.Getenv("STORAGE_NAME") // TODO put in documentation
	name = strings.ReplaceAll(name, ".", "_")
	return "storage/" + name
}

func Init() Storager {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	path := filepath.Dir(ex) + "/" + getStorageName()
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
