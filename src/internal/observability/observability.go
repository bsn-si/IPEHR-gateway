package observability

import (
	"context"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/metrics"
	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
)

type Config struct {
	ServiceName string `json:"service_name"`

	CollectMetrics bool `json:"collect_metrics"`
	MetricsPort    int  `json:"metrics_port"`

	CollectTraces  bool   `json:"collect_traces"`
	JaegerEndpoint string `json:"jaeger_endpoint"`
}

var cfg Config

func Setup(newCfg Config) {
	cfg = newCfg

	if cfg.CollectMetrics {
		metrics.SetupMetrics(cfg.ServiceName)
		metrics.RunMetricsServer(cfg.MetricsPort)
	}

	if cfg.CollectTraces {
		tracer.Setup(cfg.ServiceName, cfg.JaegerEndpoint)
	}
}

func Stop(ctx context.Context) {
	if cfg.CollectMetrics {
		metrics.Stop(ctx)
	}

	if cfg.CollectTraces {
		tracer.Stop(ctx)
	}
}
