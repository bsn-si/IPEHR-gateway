package model

// ADL
// https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/ADL1.4/operation/definition_template_adl1.4_list

type TemplateResponse struct {
	TemplateID  string `json:"template_id"`
	Concept     string `json:"concept"`
	ArchetypeID string `json:"archetype_id"`
	CreatedAt   string `json:"created_timestamp"`
}

type Template struct {
	TemplateID  string
	Version     string
	VerADL      ADLVer
	MimeType    ADLType
	Body        []byte
	Concept     string
	ArchetypeID string
	CreatedAt   string
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
