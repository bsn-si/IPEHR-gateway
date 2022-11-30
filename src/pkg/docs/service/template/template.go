package template

import (
	"context"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/errors"
	"log"
)

type Service struct {
	*service.DefaultDocumentService
	parsers map[string]ADLParser
}

type ADLParser interface {
	Version() model.VerADL
	AllowedType(string) (model.ADLTypes, error)
	Parse([]byte, string) ([]byte, error)
}

func NewService(docService *service.DefaultDocumentService, parsers ...ADLParser) *Service {
	ps := make(map[string]ADLParser)

	for _, p := range parsers {
		v := p.Version()

		_, ok := ps[v]
		if ok {
			log.Fatalf("ADL parser with ver %s already implemented", v)
		}

		ps[v] = p
	}

	return &Service{
		docService,
		ps,
	}
}

func (s *Service) Parser(version string) (ADLParser, error) {
	if len(s.parsers) == 0 {
		return nil, errors.ErrIsNotExist
	}

	p, ok := s.parsers[version]
	if !ok {
		return nil, errors.ErrIsNotExist
	}

	return p, nil
}

func (*Service) GetByID(ctx context.Context, userID string, templateID string) (*model.Template, error) {
	return nil, errors.ErrNotImplemented
}
