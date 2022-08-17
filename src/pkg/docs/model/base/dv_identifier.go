package base

// DvIdentifier
// Type for representing identifiers of real-world entities. Typical identifiers include drivers licence
// number, social security number, veterans affairs number, prescription id, order id, and so on.
// https://specifications.openehr.org/releases/RM/latest/data_types.html#_dv_identifier_class
type DvIdentifier struct {
	Issuer   string `json:"issuer,omitempty"`
	Assigner string `json:"assigner,omitempty"`
	ID       string `json:"id"`
	Type     string `json:"type,omitempty"`
}
