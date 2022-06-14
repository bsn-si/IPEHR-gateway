package access

import (
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service"
)

type AccessService struct {
	Doc *service.DefaultDocumentService
	Cfg *config.Config
}

func NewAccessService(docService *service.DefaultDocumentService, cfg *config.Config) *AccessService {
	return &AccessService{
		Doc: docService,
		Cfg: cfg,
	}
}
