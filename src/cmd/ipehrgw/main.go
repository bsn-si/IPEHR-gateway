package main

// Generating swagger doc spec
//go:generate swag init --parseDependency -g ../../pkg/api/api.go -d ../../pkg/api -o ../../pkg/api/docs

import (
	"flag"

	"hms/gateway/pkg/api"
	_ "hms/gateway/pkg/api/docs"
	"hms/gateway/pkg/config"
)

func main() {

	var (
		cfgPath = flag.String("config", "./config.json", "config file path")
	)
	flag.Parse()

	cfg := config.New(*cfgPath)
	err := cfg.Reload()
	if err != nil {
		panic(err)
	}

	a := api.New(cfg).Build()

	if err = a.Run(cfg.Host); err != nil {
		panic(err)
	}
}
