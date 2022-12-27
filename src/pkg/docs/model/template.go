package model

import "encoding/xml"

// ADL
// https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/ADL1.4/operation/definition_template_adl1.4_list

type TemplateXML struct {
	XMLName    xml.Name `xml:"template"`
	TemplateID string   `xml:"template_id>value"`
	UID        string   `xml:"uid>value"`
	Concept    string   `xml:"concept"`
	Definition struct {
		ArchetypeID string `xml:"archetype_id,attr"` // TODO maybe we should generate it?
	} `xml:"definition"`
}

type Template struct {
	TemplateID  string  `json:"template_id"`
	UID         string  `json:"-" msgpack:"-"`
	Version     string  `json:"version,omitempty" msgpack:"-"`
	VerADL      ADLVer  `json:"-" msgpack:"-"`
	MimeType    ADLType `json:"-" msgpack:"-"`
	Body        []byte  `json:"-" msgpack:"-"`
	Concept     string  `json:"concept"`
	ArchetypeID string  `json:"archetype_id,omitempty" msgpack:",omitempty"`
	CreatedAt   string  `json:"created_timestamp" msgpack:"-"`
}

// Version show as which service we should use for parsing template
type ADLVer = string

const (
	VerADL1_4 ADLVer = "adl1.4"
	VerADL2   ADLVer = "adl2"
)

type ADLType = string

const (
	ADLTypeXML  ADLType = "application/xml"
	ADLTypeJSON ADLType = "application/openehr.wt+json"
	//ADLTypeJSON ADLType = "text/plain" // for ADL2
	//ADLTypeJSON ADLType = "application/openehr.nc.flat+json"
	//ADLTypeJSON ADLType = "application/openehr.tds2+xml"
)
