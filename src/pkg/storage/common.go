package storage

import (
	"hms/gateway/pkg/storage/localfile"
	"log"
	"os"
)

var Storage Storager

func Init() Storager {
	if Storage == nil {
		cfg := localfile.Config{
			BasePath: os.Getenv("HOME") + "/Projects/bsn/hms/gateway/data/",
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
