package base

// Section
// Represents a heading in a heading structure, or section tree.
// Created according to archetyped structures for typical headings such as SOAP, physical examination,
// but also pathology result heading structures. Should not be used instead of ENTRY hierarchical structures.
// https://specifications.openehr.org/releases/RM/Release-1.0.2/ehr.html#_section_class
type Section struct {
	ContentItem
	Items *[]ContentItem `json:"items,omitempty"`
}
