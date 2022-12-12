package model

import "hms/gateway/pkg/docs/model/base"

// https://specifications.openehr.org/releases/RM/latest/common.html#_contribution_class
type Contribution struct {
	UID base.UIDBasedID `json:"uid"`
	//Versions *[]base.ObjectRef `json:"versions"`
	Versions *[]Version    `json:"versions"`
	Audit    *AuditDetails `json:"audit"`
}

//type ContributionResponse struct{}

type Version struct {
	Type           *base.ItemType   `json:"_type"`
	Contribution   base.ObjectRef   `json:"contribution"`
	CommitAudit    AuditDetails     `json:"commit_audit"`
	UID            *base.UIDBasedID `json:"uid"`
	Data           interface{}      `json:"data"` // TODO or *base.PartyProxy ???, or generics
	LifecycleState base.DvCodedText `json:"lifecycle_state"`
}

//type Version[T any] struct {
//	Type           *base.ItemType   `json:"_type"`
//	Contribution   base.ObjectRef   `json:"contribution"`
//	CommitAudit    AuditDetails     `json:"commit_audit"`
//	UID            *base.UIDBasedID `json:"uid"`
//	Data           *unknownType[T]       `json:"data"` // TODO or *base.PartyProxy ???, or generics
//	LifecycleState base.DvCodedText `json:"lifecycle_state"`
//}

// AUDIT_DETAILS
// The set of attributes required to document the committal of an information item to a repository.
// https://specifications.openehr.org/releases/RM/latest/common.html#_audit_details_class
type AuditDetails struct {
	Type         *base.ItemType    `json:"_type"`
	SystemID     string            `json:"system_id"`
	TimeCommited *base.DvDateTime  `json:"time_committed,omitempty"`
	ChangeType   *base.DvCodedText `json:"change_type,omitempty"`
	Committer    *base.PartyProxy  `json:"committer,omitempty"`
	Description  *base.DvText      `json:"description,omitempty"`
}
