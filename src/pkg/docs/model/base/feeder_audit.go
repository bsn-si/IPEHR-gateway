package base

// FeederAudit
// The FEEDER_AUDIT class defines the semantics of an audit trail which is constructed to describe the
// origin of data that have been transformed into openEHR form and committed to the system.
// https://specifications.openehr.org/releases/RM/latest/common.html#_feeder_audit_class
type FeederAudit struct {
	OriginatingSystemItemIds *[]DvIdentifier     `json:"originating_system_item_ids,omitempty"`
	FeederSystemItemIds      *[]DvIdentifier     `json:"feeder_system_item_ids,omitempty"`
	OriginalContent          DvEncapsulated      `json:"original_content"`
	OriginatingSystemAudit   FeederAuditDetails  `json:"originating_system_audit"`
	FeederSystemAudit        *FeederAuditDetails `json:"feeder_system_audit,omitempty"`
}
