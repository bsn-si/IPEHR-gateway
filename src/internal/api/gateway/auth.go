package gateway

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func auth(a *API, exceptions ...string) func(*gin.Context) {
	reCheckRegister := regexp.MustCompile("^[0-9a-f]{12}_register$")

	return func(c *gin.Context) {
		// Exception strings
		for _, ex := range exceptions {
			switch ex {
			case "userRegister":
				reqID := c.Param("reqID")

				if reCheckRegister.MatchString(reqID) {
					userID := c.Request.Header.Get("AuthUserId")
					c.Set("userID", userID)
					c.Next()
					return
				}
			}
		}

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

		c.Set("userID", userID)

		c.Next()
	}
}
