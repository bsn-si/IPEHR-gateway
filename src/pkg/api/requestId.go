package api

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/gin-gonic/gin"
)

func requestID(c *gin.Context) {
	id := make([]byte, 8)
	reqID := c.Request.Header.Get("reqId")

	if reqID == "" {
		if _, err := rand.Read(id); err != nil {
			log.Println("Make requestID error:", err)
		}

		reqID = hex.EncodeToString(id)
	}

	c.Set("reqId", reqID)
	c.Next()
}
