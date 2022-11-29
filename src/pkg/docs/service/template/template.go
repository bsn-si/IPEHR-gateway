package template

import (
	"context"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/errors"
)

type Service struct {
	*service.DefaultDocumentService
}

func NewService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		docService,
	}
}

func (*Service) GetByID(ctx context.Context, userID string, templateID string) (*model.Template, error) {
	return nil, errors.ErrNotImplemented
}
