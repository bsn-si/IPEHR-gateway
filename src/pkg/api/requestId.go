package api

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func requestId(c *gin.Context) {
	id := make([]byte, 8)
	rand.Read(id)
	c.Set("reqId", hex.EncodeToString(id))
	c.Next()
}
