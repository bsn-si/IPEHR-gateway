package base

// CodePhrase
// A fully coordinated (i.e. all coordination has been performed) term from a terminology service
// (as distinct from a particular terminology).
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_code_phrase_class
type CodePhrase struct {
	Type          ItemType `json:"_type"`
	TerminologyID ObjectID `json:"terminology_id"`
	CodeString    string   `json:"code_string,omitempty"`
	PreferredTerm string   `json:"preferred_term,omitempty"`
}
