package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/errors"
)

func auth(c *gin.Context) {
	userID := c.Request.Header.Get("AuthUserId")
	if userID == "" {
		_ = c.AbortWithError(http.StatusForbidden, errors.ErrAuthorization)
		return
	}

	c.Set("userId", userID)

	/* TODO
	signature := c.Request.Header.Get("AuthSign")
	if signature == "" {
		c.AbortWithError(http.StatusForbidden, errors.AuthorizationError)
		return
	}

	if !checkSignature(publicKey, signature) {
		c.AbortWithError(http.StatusForbidden, errors.AuthorizationError)
		return
	}
	*/

	c.Next()
}

/*
func checkSignature(pubKey, signature string) bool {
	//TODO with NaCl sign

	return true
}
*/
