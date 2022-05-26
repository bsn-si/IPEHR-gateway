package base

// DvCodedText
// A text item whose value must be the rubric from a controlled terminology, the key (i.e. the 'code') of
// which is the defining_code attribute. In other words: a DV_CODED_TEXT is a combination of a CODE_PHRASE
// (effectively a code) and the rubric of that term, from a terminology service, in the language in which
// the data were authored.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_coded_text_class
type DvCodedText struct {
	DefiningCode CodePhrase `json:"defining_code"`
	DvText
}
