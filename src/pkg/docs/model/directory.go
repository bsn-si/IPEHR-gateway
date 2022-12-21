package model

import "hms/gateway/pkg/docs/model/base"

type Directory struct {
	base.Locatable
	//base.Root
	FeederAudit base.FeederAudit   `json:"feeder_audit"`
	Folders     []Directory        `json:"folders"`
	Details     base.ItemStructure `json:"details"`
	Items       []DirectoryItem    `json:"items,omitempty"`
	//Items       base.Items         `json:"items,omitempty"`
	//Items       base.ItemStructure `json:"items"`
	//details
	//ExternalRef ExternalRef `json:"external_ref"`
}

type DirectoryItem struct {
	ID        base.UIDBasedID `json:"id"`
	Type      base.ItemType   `json:"type"`
	Namespace string          `json:"namespace"`
}

//type DirectoryDetail struct {
//}

//type ContributionVersion struct {
//	Type                base.ItemType    `json:"_type"`
//	Contribution        base.ObjectRef   `json:"contribution"`
//	CommitAudit         AuditDetails     `json:"commit_audit"`
//	UID                 base.UIDBasedID  `json:"uid"`
//	PrecedingVersionUID base.UIDBasedID  `json:"preceding_version_uid"`
//	LifecycleState      base.DvCodedText `json:"lifecycle_state"`
//	Data                base.Root        `json:"data"`
//}
