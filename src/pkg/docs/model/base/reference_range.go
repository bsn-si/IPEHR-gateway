package base

// ReferenceRange
// Defines a named range to be associated with any DV_ORDERED datum. Each such range is particular to
// the patient and context, e.g. sex, age, and any other factor which affects ranges.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_reference_range_class
type ReferenceRange[T any] struct {
	Meaning string      `json:"meaning"`
	Range   Interval[T] `json:"range"`
}
