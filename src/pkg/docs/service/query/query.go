package query

import (
	"hms/gateway/pkg/docs/service"
)

type Service struct {
	Doc *service.DefaultDocumentService
}

func NewQueryService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		Doc: docService,
	}
}
