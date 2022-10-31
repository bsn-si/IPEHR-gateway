package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/errors"
)

func auth(a *API) func(*gin.Context) {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		userID := c.Request.Header.Get("AuthUserId")

		if tokenString == "" || userID == "" {
			_ = c.AbortWithError(http.StatusForbidden, errors.ErrAuthorization)
			return
		}

		userService := a.User.service
		err := userService.VerifyAccess(userID, tokenString)

		if err != nil {
			log.Println(err)

			_ = c.AbortWithError(http.StatusForbidden, errors.ErrAuthorization)

			return
		}

		c.Set("userId", userID)

		c.Next()
	}
}
