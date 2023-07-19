package metrics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	server *http.Server
)

func RunMetricsServer(port int) {
	if server != nil {
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func Stop(ctx context.Context) {
	if server == nil {
		return
	}

	cCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(cCtx); err != nil {
		log.Fatal(err)
	}
}
