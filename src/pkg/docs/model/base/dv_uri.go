package base

// DvUri
// A reference to an object which structurally conforms to the Universal Resource Identifier (URI)
// RFC-3986 standard.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_uri_class
type DvUri struct {
	Value string `json:"value"`
}
