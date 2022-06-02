package model

import (
	"hms/gateway/pkg/docs/model/base"
)

// The EHR object is the root object and access point of an EHR for a subject of care
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_ehr_class
type EHR struct {
	SystemId      base.ObjectId     `json:"system_id"`
	EhrId         base.ObjectId     `json:"ehr_id"`
	Contributions *[]base.ObjectRef `json:"contributions,omitempty"`
	EhrStatus     base.ObjectRef    `json:"ehr_status"`
	EhrAccess     base.ObjectRef    `json:"ehr_access"`
	Compositions  *[]Composition    `json:"content,omitempty"`
	Directory     *base.ObjectRef   `json:"directory,omitempty"`
	TimeCreated   base.DvDateTime   `json:"time_created"`
	Folders       *[]base.ObjectRef `json:"folders,omitempty"`
	// "Virtual" Subject field to store external subject within EHR document. Because it is not
	// described in the official specification
	Subject struct {
		ExternalRef ExternalRef
	} `json:"-"`
}
