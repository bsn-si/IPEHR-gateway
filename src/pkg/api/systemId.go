package api

import (
	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/common"
)

func ehrSystemID(c *gin.Context) {
	ehrSystemID := c.Request.Header.Get("EhrSystemId")

	if ehrSystemID == "" {
		ehrSystemID = common.EhrSystemID
	}

	c.Set("ehrSystemID", ehrSystemID)

	c.Next()
}
