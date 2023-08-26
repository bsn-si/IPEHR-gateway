package gateway

import (
	"net/http"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func timeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(common.WebRequestTimeout),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.AbortWithStatus(http.StatusRequestTimeout)
		}),
	)
}
