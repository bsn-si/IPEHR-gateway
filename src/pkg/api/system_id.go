package api

import (
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *API) EhrSystemID(c *gin.Context) {
	ehrSystemID, err := base.NewEhrSystemID(c.Request.Header.Get("EhrSystemId"))

	if err != nil {
		_ = c.AbortWithError(http.StatusForbidden, errors.ErrIncorrectRequest)
		return
	}

	c.Set("ehrSystemID", *ehrSystemID)

	c.Next()
}
