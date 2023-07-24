package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var Middleware func(c *gin.Context) = func(c *gin.Context) {
	c.Next()
}

func middleware(c *gin.Context) {
	startTime := time.Now()
	reqSize := computeApproximateRequestSize(c.Request)

	c.Next()

	status := strconv.Itoa(c.Writer.Status())

	attrs := []attribute.KeyValue{
		attribute.String("method", c.Request.Method),
		attribute.String("status", status),
		attribute.String("url", c.Request.URL.Path),
		attribute.Int("status_code", c.Writer.Status()),
	}

	opt := metric.WithAttributes(attrs...)

	elapsed := time.Since(startTime).Milliseconds()
	respSize := int64(c.Writer.Size())

	ctx := c.Request.Context()
	requestCounter.Add(ctx, 1, opt)
	requestDuration.Record(ctx, elapsed, metric.WithAttributes(append(attrs, attribute.String("handler", c.HandlerName()))...))

	requestSize.Record(ctx, reqSize, opt)
	responseSize.Record(ctx, respSize, opt)
}

func computeApproximateRequestSize(r *http.Request) int64 {
	size := int64(0)
	if r.URL != nil {
		size = int64(len(r.URL.Path))
	}

	size += int64(len(r.Method))
	size += int64(len(r.Proto))

	for name, values := range r.Header {
		size += int64(len(name))
		for _, value := range values {
			size += int64(len(value))
		}
	}

	size += int64(len(r.Host))

	if r.ContentLength != -1 {
		size += r.ContentLength
	}

	return size
}

// func ginMetricHandle(ctx *gin.Context, start time.Time) {
// 	r := ctx.Request
// 	w := ctx.Writer

// 	// set request total
// 	_ = m.GetMetric(metricRequestTotal).Inc(nil)

// 	// set uv
// 	if clientIP := ctx.ClientIP(); !bloomFilter.Contains(clientIP) {
// 		bloomFilter.Add(clientIP)
// 		_ = m.GetMetric(metricRequestUVTotal).Inc(nil)
// 	}

// 	// set uri request total
// 	_ = m.GetMetric(metricURIRequestTotal).Inc([]string{ctx.FullPath(), r.Method, strconv.Itoa(w.Status())})

// 	// set request body size
// 	// since r.ContentLength can be negative (in some occasions) guard the operation
// 	if r.ContentLength >= 0 {
// 		_ = m.GetMetric(metricRequestBody).Add(nil, float64(r.ContentLength))
// 	}

// 	// set slow request
// 	latency := time.Since(start)
// 	if int32(latency.Seconds()) > m.slowTime {
// 		_ = m.GetMetric(metricSlowRequest).Inc([]string{ctx.FullPath(), r.Method, strconv.Itoa(w.Status())})
// 	}

// 	// set request duration
// 	_ = m.GetMetric(metricRequestDuration).Observe([]string{ctx.FullPath()}, latency.Seconds())

// 	// set response size
// 	if w.Size() > 0 {
// 		_ = m.GetMetric(metricResponseBody).Add(nil, float64(w.Size()))
// 	}
// }
