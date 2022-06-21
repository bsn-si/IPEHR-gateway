package base

// TermMapping
// Represents a coded term mapped to a DV_TEXT, and the relative match of the target term with respect
// to the mapped item.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_term_mapping_class
type TermMapping struct {
	Match   rune         `json:"match"`
	Purpose *DvCodedText `json:"purpose,omitempty"`
	Target  CodePhrase   `json:"target"`
}
