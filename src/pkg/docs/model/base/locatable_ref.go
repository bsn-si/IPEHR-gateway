package base

// LocatableRef
// Purpose Reference to a LOCATABLE instance inside the top-level content structure inside a VERSION<T>;
// the path attribute is applied to the object that VERSION.data points to.
// https://specifications.openehr.org/releases/RM/Release-1.0.2/support.html#_locatable_ref_class
type LocatableRef struct {
	Path string     `json:"path,omitempty"`
	ID   UIDBasedID `json:"id"`
	ObjectRef
}
