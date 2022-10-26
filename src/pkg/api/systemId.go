package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/errors"
)

func ehrSystemID(c *gin.Context) {
	ehrSystemID := c.Request.Header.Get("EhrSystemId")
	if ehrSystemID == "" {
		_ = c.AbortWithError(http.StatusForbidden, errors.ErrIncorrectRequest)
		return
	}

	c.Set("ehrSystemID", ehrSystemID)

	c.Next()
}
