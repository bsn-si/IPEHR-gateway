package api

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"

	log "hms/gateway/pkg/logging"
)

func requestID(c *gin.Context) {
	id := make([]byte, 6)
	reqID := c.Request.Header.Get("reqID")

	if reqID == "" {
		if _, err := rand.Read(id); err != nil {
			log.Println("Make requestID error:", err)
		}

		reqID = hex.EncodeToString(id)
	}

	c.Set("reqID", reqID)
	c.Header("RequestId", reqID)
	c.Next()
}
