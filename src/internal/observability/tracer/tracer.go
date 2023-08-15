package tracer

import (
	"context"
	"log"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var tracerProvider *sdktrace.TracerProvider
var svcName string

func GetTracer(opts ...trace.TracerOption) trace.Tracer {
	return otel.Tracer(svcName, opts...)
}

func Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return GetTracer().Start(ctx, name, opts...)
}

func Setup(serviceName, url string) {
	if tracerProvider != nil {
		return
	}

	svcName = serviceName

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		log.Fatal(err)
	}

	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(getServiceResource(serviceName)),
	)

	otel.SetTracerProvider(tracerProvider)

	Middleware = otelgin.Middleware(serviceName,
		otelgin.WithTracerProvider(tracerProvider),
	)
}

func Stop(ctx context.Context) {
	if tracerProvider == nil {
		return
	}

	if err := tracerProvider.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
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
