package api

import (
	"github.com/gin-gonic/gin"
)

func ehrSystemID(c *gin.Context) {
	ehrSystemID := c.Request.Header.Get("EhrSystemId")
	if ehrSystemID != "" {
		c.Set("ehrSystemID", ehrSystemID)
	}

	c.Next()
}
