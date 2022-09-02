package base

// CareEntry
// The abstract parent of all clinical ENTRY subtypes.
// A CARE_ENTRY defines protocol and guideline attributes for all clinical Entry subtypes.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_care_entry_class
type CareEntry struct {
	Protocol    ItemStructure `json:"protocol,omitempty"`
	GuidelineID ObjectRef     `json:"guideline_id"`
	Entry
}
