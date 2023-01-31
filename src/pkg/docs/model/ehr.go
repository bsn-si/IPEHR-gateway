package model

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
)

// EHR model info
// The EHR object is the root object and access point of an EHR for a subject of care
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_ehr_class
type EHR struct {
	SystemID      base.HierObjectID `json:"system_id"`
	EhrID         base.HierObjectID `json:"ehr_id"`
	Contributions []base.ObjectRef  `json:"contributions,omitempty"`
	EhrStatus     base.ObjectRef    `json:"ehr_status"`
	EhrAccess     base.ObjectRef    `json:"ehr_access"`
	Compositions  []*Composition    `json:"compositions,omitempty"`
	Directory     *base.ObjectRef   `json:"directory,omitempty"`
	TimeCreated   base.DvDateTime   `json:"time_created"`
	Folders       []base.ObjectRef  `json:"folders,omitempty"`
}
