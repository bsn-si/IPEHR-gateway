package adl14

import (
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (*Parser) Version() model.VerADL {
	return model.VerADL1_4
}

func (*Parser) AllowedType(s string) (model.ADLTypes, error) {
	switch s {
	case model.ADLTypeXML:
		return model.ADLTypeXML, nil
	case model.ADLTypeJSON:
		return model.ADLTypeXML, nil
	}

	return "", errors.ErrNotFound
}

func (*Parser) Parse(b []byte, t model.ADLTypes) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}
