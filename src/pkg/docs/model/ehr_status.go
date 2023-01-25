package model

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
)

// Single object per EHR containing various EHR-wide status flags and settings,
// including whether this EHR can be queried, modified etc.
// This object is always modifiable, in order to change the status of the EHR as a whole.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_ehr_status_class
type EhrStatus struct {
	base.Locatable
	Subject      base.PartySelf      `json:"subject"`
	IsQueryable  bool                `json:"is_queryable"`
	IsModifable  bool                `json:"is_modifiable"`
	OtherDetails *base.ItemStructure `json:"other_details,omitempty"`
}
