package tracer

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/semconv/v1.17.0/httpconv"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var Middleware = func(c *gin.Context) {
	c.Next()
}

func middleware(c *gin.Context) {
	savedCtx := c.Request.Context()
	defer func() {
		c.Request = c.Request.WithContext(savedCtx)
	}()

	opts := []oteltrace.SpanStartOption{
		oteltrace.WithAttributes(httpconv.ServerRequest(svcName, c.Request)...),
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
	}

	spanName := c.FullPath()
	if spanName == "" {
		spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
	} else {
		rAttr := semconv.HTTPRouteKey.String(spanName)
		opts = append(opts, oteltrace.WithAttributes(rAttr))
	}

	ctx, span := GetTracer().Start(savedCtx, spanName, opts...)
	defer span.End()

	// pass the span through the request context
	c.Request = c.Request.WithContext(ctx)

	// serve the request to the next middleware
	c.Next()

	status := c.Writer.Status()
	span.SetStatus(httpconv.ServerStatus(status))

	if status != 0 {
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(status))
	}

	if len(c.Errors) > 0 {
		span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
	}
}
