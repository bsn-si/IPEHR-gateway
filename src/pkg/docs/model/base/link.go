package base

// Link
// The LINK type defines a logical relationship between two items, such as two ENTRYs or an ENTRY and
// a COMPOSITION. Links can be used across compositions, and across EHRs. Links can potentially be used
// between interior (i.e. non archetype root) nodes, although this probably should be prevented in
// archetypes. Multiple LINKs can be attached to the root object of any archetyped structure to give
// the effect of a 1→N link.
// https://specifications.openehr.org/releases/RM/latest/common.html#_link_class
type Link struct {
	Meaning DvText `json:"meaning"`
	Type    DvText `json:"type"`
	Target  DvURI  `json:"target"`
}
