package logger

import (
	"math"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func HttpMiddleware(logger *ServiceLogger) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()

		c.Next()

		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		fields := Fields{
			"path":       path,
			"method":     c.Request.Method,
			"requestId":  c.GetString("reqID"),
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
			"headers":    c.Request.Header,
		}

		switch {
		case statusCode >= 500:
			logger.WithFields(fields).Errorln("[GIN]", fields["method"], fields["path"])
		case statusCode >= 400:
			logger.WithFields(fields).Warnln("[GIN]", fields["method"], fields["path"])
		default:
			logger.WithFields(fields).Debugln("[GIN]", fields["method"], fields["path"])
		}
	}
}
