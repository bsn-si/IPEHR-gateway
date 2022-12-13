package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/errors"
)

func ehrSystemID(c *gin.Context) {
	ehrSystemID := c.Request.Header.Get("EhrSystemId")
	if ehrSystemID == "" {
		_ = c.AbortWithError(http.StatusBadRequest, errors.ErrIncorrectRequest)
		return
	}

	c.Set("ehrSystemID", ehrSystemID)

	c.Next()
}
