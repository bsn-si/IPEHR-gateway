package base

// FeederAuditDetails
// Audit details for any system in a feeder system chain. Audit details here means the general notion of
// who/where/when the information item to which the audit is attached was created.
// https://specifications.openehr.org/releases/RM/latest/common.html#_feeder_audit_details_class
type FeederAuditDetails struct {
	SystemID     string           `json:"system_id"`
	Location     *PartyIdentified `json:"location,omitempty"`
	Subject      *PartyProxy      `json:"subject,omitempty"`
	Provider     *PartyIdentified `json:"provider,omitempty"`
	Time         *DvDateTime      `json:"time,omitempty"`
	VersionID    string           `json:"version_id,omitempty"`
	OtherDetails Locatable        `json:"other_details"`
}
