package query

import (
	"hms/gateway/pkg/docs/service"
)

type QueryService struct {
	Doc *service.DefaultDocumentService
}

func NewQueryService(docService *service.DefaultDocumentService) *QueryService {
	return &QueryService{
		Doc: docService,
	}
}
