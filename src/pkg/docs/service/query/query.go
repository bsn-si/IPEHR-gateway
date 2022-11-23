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

func (*Service) Validate(data []byte) bool {
	return false
}

func (*Service) Store(ctx context.Context, userID string, qType string, qualifiedQueryName string, q []byte) (model.StoredQuery, error) {
	return model.StoredQuery{}, nil
}
