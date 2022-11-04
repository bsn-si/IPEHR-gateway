package api

import (
	"crypto/rand"
	"encoding/hex"
	"hms/gateway/pkg/common"
	"log"
	"path"

	"github.com/gin-gonic/gin"
)

func requestID(c *gin.Context) {
	id := make([]byte, 6)
	reqID := c.Request.Header.Get("reqId")

	if reqID == "" {
		if _, err := rand.Read(id); err != nil {
			log.Println("Make requestID error:", err)
		}

		lastPart := path.Base(c.Request.URL.RequestURI())
		reqID = hex.EncodeToString(id) + common.RequestIDSeparator + lastPart
	}

	c.Set("reqId", reqID)
	c.Header("RequestId", reqID)
	c.Next()
}

func requestIDFromParam(c *gin.Context) {
	reqID := c.Param("reqId")

	c.Set("reqId", reqID)
	c.Header("RequestId", reqID)
	c.Next()
}
