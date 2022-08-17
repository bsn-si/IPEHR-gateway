package request

import (
	"hms/gateway/pkg/docs/service"
)

type Service struct {
	Doc *service.DefaultDocumentService
}

func NewRequestService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		Doc: docService,
	}
}
