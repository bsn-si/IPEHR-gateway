package request

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
)

type Service struct {
	Doc *service.DefaultDocumentService
}

func NewRequestService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		Doc: docService,
	}
}
