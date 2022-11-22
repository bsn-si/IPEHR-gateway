package query

import (
	"context"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
)

type Service struct {
	DefaultDocumentService *service.DefaultDocumentService
}

func NewService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		DefaultDocumentService: docService,
	}
}

func (*Service) Get(ctx context.Context, userID string, qualifiedQueryName string) ([]model.StoredQuery, error) {
	return nil, nil
}
