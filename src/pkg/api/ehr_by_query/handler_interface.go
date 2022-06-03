package ehr_by_query

import "github.com/gin-gonic/gin"

type HandlerInterface interface {
	CanProcess(c *gin.Context) bool
	Handle(c *gin.Context) ([]byte, error)
}
