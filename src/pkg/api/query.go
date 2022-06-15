package api

import (
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/query"

	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
	*query.QueryService
}

func NewQueryHandler(docService *service.DefaultDocumentService, cfg *config.Config) *QueryHandler {
	return &QueryHandler{
		query.NewQueryService(docService, cfg),
	}
}

func (h QueryHandler) ExecPost(c *gin.Context) {
}
