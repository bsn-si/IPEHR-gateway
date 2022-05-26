package base

// DvAmount
// Abstract class defining the concept of relative quantified 'amounts'. For relative quantities,
// the + and - operators are defined (unlike descendants of DV_ABSOLUTE_QUANTITY, such as the date/time
// types).
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_amount_class
type DvAmount struct {
	AccuracyIsPercent bool    `json:"accuracy_is_percent,omitempty"`
	Accuracy          float32 `json:"accuracy,omitempty"`
	DvQuantified
}
