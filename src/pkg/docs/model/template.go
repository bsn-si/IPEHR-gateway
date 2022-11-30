package model

// ADL
// https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/ADL1.4/operation/definition_template_adl1.4_list

type Template struct {
	TemplateID  string `json:"template_id"`
	Version     string `json:"version"`
	Concept     string `json:"concept"`
	ArchetypeID string `json:"archetype_id"`
	CreatedAt   string `json:"created_timestamp"`
	verADL      string
	content     string
}

// Version show as which service we should use for parsing template
type verADL = string

const (
	verADL1_4 verADL = "adl1.4"
	verADL2   verADL = "adl2"
)
