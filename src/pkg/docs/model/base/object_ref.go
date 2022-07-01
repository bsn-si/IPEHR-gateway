package base

// ObjectRef
// Class describing a reference to another object, which may exist locally or be maintained outside
// the current namespace, e.g. in another service.
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_object_ref_class
type ObjectRef struct {
	ID        ObjectID `json:"id"`
	Namespace string   `json:"namespace"`
	Type      string   `json:"type"`
}
