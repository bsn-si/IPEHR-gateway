package api

import (
	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/access"
)

type AccessHandler struct {
	*access.AccessService
}

func NewAccessHandler(docService *service.DefaultDocumentService, cfg *config.Config) *AccessHandler {
	return &AccessHandler{
		access.NewAccessService(docService, cfg),
	}
}

func (h AccessHandler) GroupCreate(c *gin.Context) {
	//
}
