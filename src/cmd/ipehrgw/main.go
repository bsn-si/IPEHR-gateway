package main

// Generating swagger doc spec//
//go:generate swag fmt -g ../../pkg/api/api.go
//go:generate swag init --parseDepth 1 -g ./../../pkg/api/api.go -d ./../../pkg/api,./../../pkg/docs/model,./../../pkg/docs/service/processing,./../../pkg/user/model -o ./../../pkg/api/docs

import (
	"flag"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/api"
	_ "github.com/bsn-si/IPEHR-gateway/src/pkg/api/docs"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
)

func main() {
	var (
		cfgPath = flag.String("config", "./config.json", "config file path")
	)

	flag.Parse()

	cfg, err := config.New(*cfgPath)
	if err != nil {
		panic(err)
	}

	infra := infrastructure.New(cfg)

	a := api.New(cfg, infra).Build()
	if err = a.Run(cfg.Host); err != nil {
		panic(err)
	}
}
