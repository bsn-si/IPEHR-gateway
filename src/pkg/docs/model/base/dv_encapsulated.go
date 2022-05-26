package base

// DvEncapsulated
// Abstract class defining the common meta-data of all types of encapsulated data.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_encapsulated_class
type DvEncapsulated struct {
	Charset  *CodePhrase `json:"charset,omitempty"`
	Language *CodePhrase `json:"language,omitempty"`
}
