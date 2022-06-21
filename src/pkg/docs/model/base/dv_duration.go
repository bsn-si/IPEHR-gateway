package base

// DvDuration
// Represents a period of time with respect to a notional point in time, which is not specified.
// A sign may be used to indicate the duration is backwards in time rather than forwards.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_duration_class
type DvDuration struct {
	Value string `json:"value"`
	DvAmount
}
