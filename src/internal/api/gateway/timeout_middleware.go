package gateway

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func timeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(common.QueryExecutionTimeout),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
	)
}
