package metrics

import (
	"log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var (
	requestCounter  instrument.Int64Counter
	requestDuration instrument.Int64Histogram
	requestSize     instrument.Int64Histogram
	responseSize    instrument.Int64Histogram
)

func SetupMetrics(serviceName string) {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(getServiceResource(serviceName)),
	)
	global.SetMeterProvider(provider)
	meter := global.Meter(serviceName)

	requestDuration, err = meter.Int64Histogram("request_duration",
		instrument.WithDescription("Request duration in milliseconds"),
		instrument.WithUnit(unit.Milliseconds),
	)
	if err != nil {
		log.Fatalf("error on create request_duration histogram: %v", err)
	}

	requestSize, err = meter.Int64Histogram("request_size",
		instrument.WithDescription("Request size"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		log.Fatalf("error on create request_size histogram: %v", err)
	}

	responseSize, err = meter.Int64Histogram("response_size",
		instrument.WithDescription("Response size"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		log.Fatalf("error on create response_size histogram: %v", err)
	}

	requestCounter, err = meter.Int64Counter("request_counter",
		instrument.WithDescription("Total requests count"),
	)
	if err != nil {
		log.Fatalf("error on request counter: %v", err)
	}

	if err != nil {
		log.Fatalf("error on blocks counter: %v", err)
	}

	Middleware = middleware
}

func getServiceResource(serviceName string) *resource.Resource {
	defaultOpts := resource.Default()
	attrOpts := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String("v0.1.0"),
		attribute.String("environment", "production"),
	)

	r, err := resource.Merge(defaultOpts, attrOpts)

	if err != nil {
		log.Fatal(err)
	}

	return r
}
