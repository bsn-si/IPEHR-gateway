package logging

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Middleware(logger *ServiceLogger) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()

		fields := Fields{
			"path":      path,
			"method":    c.Request.Method,
			"requestId": c.GetString("reqID"),
			"hostname":  hostname,
			"clientIP":  c.Request.Referer(),
			"referer":   c.Request.Referer(),
			"userAgent": c.Request.UserAgent(),
			"headers":   c.Request.Header,
		}
		ctx := ContextWithFields(c.Request.Context(), fields)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		statusCode := c.Writer.Status()

		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		fields["statusCode"] = statusCode
		fields["dataLength"] = dataLength
		fields["latency"] = time.Since(start).Seconds()

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
