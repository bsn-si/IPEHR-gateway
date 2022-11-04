package api

import (
	"hms/gateway/pkg/common"
	"log"
	"net/http"
	"strings"

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

func authWithExceptions(a *API, exceptPath string) func(*gin.Context) {
	return func(c *gin.Context) {
		requestID := c.GetString("reqId")
		part := ""

		strArr := strings.Split(requestID, common.RequestIDSeparator)
		if len(strArr) > 1 {
			part = strArr[len(strArr)-1]
		}

		if part != exceptPath {
			//c.HandlerName()
			auth(a)(c)
		} else {
			userID := c.Request.Header.Get("AuthUserId")
			c.Set("userId", userID)
		}

		c.Next()
	}
}
