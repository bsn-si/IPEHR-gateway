package main

// Generating swagger doc spec
//go:generate swag init --parseDependency -g ../../pkg/api/api.go -d ../../pkg/api -o ../../pkg/api/docs

import (
	"hms/gateway/pkg/api"
	_ "hms/gateway/pkg/api/docs"
	"hms/gateway/pkg/config"
)

func main() {

	cfg := config.GetConfig("../config.json")

	a := api.New()

	r := a.Build()

	err := r.Run(cfg.Host)
	if err != nil {
		panic(err)
	}
}
