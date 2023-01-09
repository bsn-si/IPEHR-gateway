package main

// Generating swagger doc spec//
//go:generate swag fmt -g ../../pkg/api/api.go
//go:generate swag init --parseDependency -g ../../pkg/api/api.go -d ../../pkg/api -o ../../pkg/api/docs

import (
	"flag"
	"log"

	"hms/gateway/pkg/api"
	_ "hms/gateway/pkg/api/docs"
	_ "hms/gateway/pkg/aqlquerier"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service/query"
	"hms/gateway/pkg/infrastructure"

	"github.com/jmoiron/sqlx"
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

	conn, err := sqlx.Open("aql", "")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	executer := query.NewQueryExecuterService(conn)

	a := api.New(cfg, infra, executer).Build()

	if err = a.Run(cfg.Host); err != nil {
		panic(err)
	}
}
