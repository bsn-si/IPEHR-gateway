package base

// Interval abstraction, featuring upper and lower limits that may be open or closed, included or not
// included. Interval of ordered items.
// https://specifications.openehr.org/releases/BASE/latest/foundation_types.html#_interval_class
type Interval[T any] struct {
	Lower          T    `json:"lower,omitempty"`
	Upper          T    `json:"upper,omitempty"`
	LowerUnbounded bool `json:"lower_unbounded"`
	UpperUnbounded bool `json:"upper_unbounded"`
	LowerIncluded  bool `json:"lower_included"`
	UpperIncluded  bool `json:"upper_included"`
}
