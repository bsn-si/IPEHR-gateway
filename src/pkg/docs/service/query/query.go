package query

import (
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service"
)

type QueryService struct {
	Doc *service.DefaultDocumentService
	Cfg *config.Config
}

func NewQueryService(docService *service.DefaultDocumentService, cfg *config.Config) *QueryService {
	return &QueryService{
		Doc: docService,
		Cfg: cfg,
	}
}
