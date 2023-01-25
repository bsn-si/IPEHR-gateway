package model

import "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"

// EventContext
// Documents the context information of a healthcare event involving the subject of care and the health
// system. The context information recorded here are independent of the attributes recorded in the version
// audit, which document the system interaction context, i.e. the context of a user interacting with
// the health record system.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_event_context_class
type EventContext struct {
	StartTime          base.DvDateTime       `json:"start_time"`
	EndTime            *base.DvDateTime      `json:"end_time,omitempty"`
	Location           *string               `json:"location,omitempty"`
	Setting            base.DvCodedText      `json:"setting"`
	OtherContext       *base.ItemStructure   `json:"other_context,omitempty"`
	HealthCareFacility *base.PartyIdentified `json:"health_care_facility,omitempty"`
	Participations     []base.Participation  `json:"participations,omitempty"`
}
