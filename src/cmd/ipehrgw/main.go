package main

// Generating swagger doc spec//
//go:generate swag fmt -g ../../internal/api/gateway/api.go
//go:generate swag init --parseDepth 1 -g ./../../api/gateway/api.go -d ./../../internal/api/gateway,./../../pkg/docs/model,./../../pkg/docs/service/processing,./../../pkg/user/model -o ./../../internal/api/gateway/docs

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/bsn-si/IPEHR-gateway/src/internal/api/gateway/docs"
	"github.com/bsn-si/IPEHR-gateway/src/internal/observability"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
	"github.com/circonus-labs/circonus-gometrics/api"
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	observability.Setup(cfg.Observability)

	infra := infrastructure.New(cfg)

	handler := api.New(cfg, infra).Build()

	server := http.Server{
		Addr:    cfg.Host,
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	observability.Stop(ctx)

	log.Println("Server exiting")
}
