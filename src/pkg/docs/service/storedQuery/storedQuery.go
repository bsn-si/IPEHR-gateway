package storedQuery

import (
	"context"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
)

type Service struct {
	*service.DefaultDocumentService
	defaultGroupAccess *model.GroupAccess
}

func NewService(docService *service.DefaultDocumentService) *Service {
	s := &Service{
		DefaultDocumentService: docService,
	}

	return s
}

func (*Service) Get(ctx context.Context, userID string, qualifiedQueryName string) ([]*model.StoredQuery, error) {
	return nil, nil
}
