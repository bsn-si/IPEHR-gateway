package base

// DvQuantified
// Abstract class defining the concept of true quantified values, i.e. values which are not only ordered,
// but which have a precise magnitude.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_quantified_class
type DvQuantified struct {
	MagnitudeStatus bool        `json:"magnitude_status,omitempty"`
	Accuracy        interface{} `json:"accuracy,omitempty"`
	DvOrdered
}
