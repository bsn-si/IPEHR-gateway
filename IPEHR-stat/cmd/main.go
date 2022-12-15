package main

import (
	"flag"

	"ipehr/stat/pkg/api"
	"ipehr/stat/pkg/config"
	"ipehr/stat/pkg/infrastructure"
	"ipehr/stat/pkg/service/syncer"

	"github.com/gin-contrib/cors"
)

func main() {
	cfgPath := flag.String("config", "./config.json", "config file path")

	flag.Parse()

	cfg := config.New(*cfgPath)

	infra := infrastructure.New(cfg)
	defer infra.Close()

	syncer.New(
		infra.DB,
		infra.EthClient,
		syncer.Config(cfg.Sync),
	).Start()

	a := api.New(cfg, infra).Build()

	//TODO complete CORS config
	a.Use(cors.Default())

	err := a.Run(cfg.Host)
	if err != nil {
		panic(err)
	}
}
