package main

// Generating swagger doc spec//
//go:generate swag fmt -g ../../internal/api/gateway/api.go
//go:generate swag init --parseDepth 1 -g ./../../api/gateway/api.go -d ./../../internal/api/gateway,./../../pkg/docs/model,./../../pkg/docs/service/processing,./../../pkg/user/model -o ./../../internal/api/gateway/docs

import (
	"flag"

	"github.com/bsn-si/IPEHR-gateway/src/internal/api/gateway"
	_ "github.com/bsn-si/IPEHR-gateway/src/internal/api/gateway/docs"
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

	a := gateway.New(cfg, infra).Build()
	if err = a.Run(cfg.Host); err != nil {
		panic(err)
	}
}
