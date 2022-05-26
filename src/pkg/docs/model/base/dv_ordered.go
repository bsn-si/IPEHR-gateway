package base

// DvOrdered
// Abstract class defining the concept of ordered values, which includes ordinals as well as true
// quantities.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_ordered_class
type DvOrdered struct {
	NormalStatus         *CodePhrase       `json:"normal_status,omitempty"`
	NormalRange          *Interval         `json:"normal_range,omitempty"`
	OtherReferenceRanges *[]ReferenceRange `json:"other_reference_ranges,omitempty"`
}
