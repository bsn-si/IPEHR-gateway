package adl14

import (
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
)

type Parser struct{}

func NewADLParser() *Parser {
	return &Parser{}
}

func (*Parser) Version() model.VerADL {
	return model.VerADL1_4
}

func (*Parser) Parse(string) (interface{}, error) {
	return nil, errors.ErrNotImplemented
}
