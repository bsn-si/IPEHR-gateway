package adl2

import (
	"hms/gateway/pkg/docs/model"
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
