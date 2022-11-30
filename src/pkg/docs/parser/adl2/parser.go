package adl2

import (
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (*Parser) AllowedType(s string) (model.ADLTypes, error) {
	return "", errors.ErrNotImplemented
}

func (*Parser) Version() model.VerADL {
	return model.VerADL2
}

func (*Parser) Parse(b []byte, t model.ADLTypes) ([]byte, error) {
	return nil, errors.ErrNotImplemented
}
