package ehr_by_query

import "github.com/gin-gonic/gin"

// GetEhrByQuery Find EHR by query parameters
func GetEhrByQuery(c *gin.Context) ([]byte, error) {
	handler, err := Resolve(c)
	if err != nil {
		return nil, err
	}

	return handler.Handle(c)
}
