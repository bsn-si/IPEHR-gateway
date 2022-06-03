package ehr_by_query

import (
	"errors"
	"github.com/gin-gonic/gin"
	"hms/gateway/pkg/api/ehr_by_query/handler"
)

// List of exists handler to process query requests
var handlers []HandlerInterface

func init() {
	handlers = []HandlerInterface{
		&handler.Subject{},
	}
}

// Resolve which handler can process query request
func Resolve(c *gin.Context) (HandlerInterface, error) {
	for _, currentHandler := range handlers {
		if currentHandler.CanProcess(c) {
			return currentHandler, nil
		}
	}
	return nil, errors.New("can't handle query")
}
