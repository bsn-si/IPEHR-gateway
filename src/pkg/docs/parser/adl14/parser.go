package adl14

import (
	"encoding/xml"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
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

func (*Parser) Validate(b []byte, mime model.ADLType) bool {
	if mime == model.ADLTypeXML {
		return xml.Unmarshal(b, new(interface{})) == nil
	}

	return false
}

func (*Parser) Parse(b []byte, mime model.ADLType) (*model.Template, error) {
	if mime == model.ADLTypeXML {
		m := model.TemplateXML{}

		err := xml.Unmarshal(b, &m)
		if err != nil {
			return nil, errors.Wrap(err, "cannot unmarshal XML type")
		}

		return &model.Template{
			TemplateID:  m.TemplateID,
			UID:         m.UID,
			ArchetypeID: m.Definition.ArchetypeID,
			VerADL:      model.VerADL1_4,
			MimeType:    mime,
			Body:        b,
			Concept:     m.Concept,
		}, nil
	}

	return nil, errors.Wrap(errors.ErrNotImplemented, "type have not implementation")
}

func (p *Parser) ParseWithFill(b []byte, mime model.ADLType) (*model.Template, error) {
	m, err := p.Parse(b, mime)
	if err != nil {
		return nil, err
	}

	m.CreatedAt = time.Now().Format(common.OpenEhrTimeFormat)
	m.Version = "1"
	m.Body = b

	return m, nil
}
