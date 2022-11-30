package model

// ADL
// https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/ADL1.4/operation/definition_template_adl1.4_list

type Template struct {
	TemplateID  string `json:"template_id"`
	Version     string `json:"version"`
	Concept     string `json:"concept"`
	ArchetypeID string `json:"archetype_id"`
	CreatedAt   string `json:"created_timestamp"`
	VerADL      VerADL
	Content     []byte
}

// Version show as which service we should use for parsing template
type VerADL = string

const (
	VerADL1_4 VerADL = "adl1.4"
	VerADL2   VerADL = "adl2"
)

type ADLTypes = string

const (
	ADLTypeXML  ADLTypes = "application/xml"
	ADLTypeJSON ADLTypes = "application/openehr.wt+json"
	//ADLTypeJSON ADLTypes = "text/plain" // for ADL2
	//ADLTypeJSON ADLTypes = "application/openehr.nc.flat+json"
	//ADLTypeJSON ADLTypes = "application/openehr.tds2+xml"
)
