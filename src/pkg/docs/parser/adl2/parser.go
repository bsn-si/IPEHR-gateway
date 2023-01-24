package adl2

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (*Parser) IsTypeAllowed(t model.ADLType) bool {
	// TODO
	return false
}

func (*Parser) Version() model.ADLVer {
	return model.VerADL2
}

func (*Parser) Validate(b []byte, t model.ADLType) bool {
	// TODO
	return false
}

func (*Parser) Parse([]byte, model.ADLType) (*model.Template, error) {
	return nil, errors.ErrNotImplemented
}

func (*Parser) ParseWithFill([]byte, model.ADLType) (*model.Template, error) {
	return nil, errors.ErrNotImplemented
}
