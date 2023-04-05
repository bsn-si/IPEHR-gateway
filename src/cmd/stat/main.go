package main

// Generating swagger doc spec//
//go:generate swag fmt -g ../internal/api/stat/api.go
//go:generate swag init --parseDependency -g ../cmd/stat/main.go -o ../internal/api/stat/docs

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/internal/api/stat"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/service/syncer"

	"github.com/gin-contrib/cors"
)

// @title        IPEHR Stat API
// @version      0.1
// @description  IPEHR Stat is an open API service for providing public statistics from the IPEHR system.

// @contact.name   API Support
// @contact.url    https://bsn.si/blockchain
// @contact.email  support@bsn.si

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      stat.ipehr.org
// host      localhost:8080
// @BasePath  /

func main() {
	cfgPath := flag.String("config", "./config.json", "config file path")

	flag.Parse()

	cfg := config.NewStatConfig(*cfgPath)

	infra := infrastructure.NewStatInfra(cfg)
	defer infra.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Interrupt)
	defer cancel()

	syncer.New(
		infra.StatsRepo,
		infra.ChunkRepo,
		infra.EthClient,
		cfg.Sync,
	).Start(ctx)

	router := stat.New(cfg, infra).Build()

	//TODO complete CORS config
	router.Use(cors.Default())

	srv := http.Server{ //nolint
		Addr:    cfg.Host,
		Handler: router,
	}

	go func() {
		log.Printf("Server start listening on host: %v", srv.Addr)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Listen server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Server shutdowning...")

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()

	if err := srv.Shutdown(stopCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
