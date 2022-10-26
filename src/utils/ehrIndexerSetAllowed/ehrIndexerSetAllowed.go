package main

import (
	"context"
	"flag"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/indexer"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
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

	ehtClient, err := ethclient.Dial(cfg.Contract.Endpoint)
	if err != nil {
		log.Fatal(err)
	}

	index := indexer.New(cfg.Contract.Address, cfg.Contract.PrivKeyPath, ehtClient)

	address := "<address>"

	_, err = index.SetAllowed(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}
}
