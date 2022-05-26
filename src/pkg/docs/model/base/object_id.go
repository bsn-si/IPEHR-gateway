package base

// ObjectId
// Ancestor class of identifiers of informational objects. Ids may be completely meaningless, in which
// case their only job is to refer to something, or may carry some information to do with the identified
// object.
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_object_id_class
type ObjectId struct {
	Type  string `json:"_type,omitempty"`
	Value string `json:"value"`
}
