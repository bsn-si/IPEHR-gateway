package api

import (
	"hms/gateway/pkg/config"

	"github.com/gin-gonic/gin"
)

func baseURL(c *config.Config) func(*gin.Context) {
	baseURL := ""
	if c != nil {
		baseURL = c.BaseURL
	}
	return func(c *gin.Context) {
		c.Set("baseURL", baseURL)
		c.Next()
	}
}
