package api

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/gin-gonic/gin"
)

func requestID(c *gin.Context) {
	id := make([]byte, 8)

	if _, err := rand.Read(id); err != nil {
		log.Println("Make requestID error:", err)
	}

	c.Set("reqId", hex.EncodeToString(id))
	c.Next()
}
