package adl14

import (
	"hms/gateway/pkg/docs/model"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (*Parser) Version() model.ADLVer {
	return model.VerADL1_4
}

func (*Parser) IsTypeAllowed(t model.ADLType) bool {
	return t == model.ADLTypeXML || t == model.ADLTypeJSON
}

func (*Parser) Validate(b []byte, t model.ADLType) bool {
	// TODO
	return false
}
