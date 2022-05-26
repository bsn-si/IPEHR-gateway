package base

// DvDateTime
// Represents an absolute point in time, specified to the second. Semantics defined by ISO 8601.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_date_time_class
type DvDateTime struct {
	Value    string      `json:"value"`
	Accuracy *DvDuration `json:"accuracy,omitempty"`
}
