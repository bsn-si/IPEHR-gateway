package base

// DvText
// A text item, which may contain any amount of legal characters arranged as e.g. words, sentences etc
// (i.e. one DV_TEXT may be more than one word). Visual formatting and hyperlinks may be included via
// markdown.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_text_class
type DvText struct {
	Value      string         `json:"value"`
	Hyperlink  *DvUri         `json:"hyperlink,omitempty"`
	Formatting string         `json:"formatting,omitempty"`
	Mappings   *[]TermMapping `json:"mappings,omitempty"`
	Language   *CodePhrase    `json:"language,omitempty"`
	Encoding   *CodePhrase    `json:"encoding,omitempty"`
}
