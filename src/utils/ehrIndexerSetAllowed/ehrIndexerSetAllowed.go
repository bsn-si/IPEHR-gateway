package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer"
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

	index := indexer.New(
		cfg.Contract.AddressEhrIndex,
		cfg.Contract.AddressAccessStore,
		cfg.Contract.AddressUsers,
		cfg.Contract.AddressDataStore,
		cfg.Contract.PrivKeyPath,
		ehtClient,
		cfg.Contract.GasTipCap,
	)

	key, err := os.ReadFile(cfg.Contract.PrivKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimSpace(string(key)))
	if err != nil {
		log.Fatal(err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	txHash, err := index.SetAllowed(context.Background(), address.String())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("txHash: ", txHash)
}
