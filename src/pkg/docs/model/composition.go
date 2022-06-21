package model

import "hms/gateway/pkg/docs/model/base"

// Composition Content of one version in a VERSIONED_COMPOSITION. A Composition is considered the unit
// of modification of the record, the unit of transmission in record Extracts, and the unit of
// attestation by authorising clinicians. In this latter sense, it may be considered equivalent to a
// signed document.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_composition_class
type Composition struct {
	Language  base.CodePhrase     `json:"language"`
	Territory base.CodePhrase     `json:"territory"`
	Category  base.DvCodedText    `json:"category"`
	Context   *EventContext       `json:"context,omitempty"`
	Composer  base.PartyProxy     `json:"composer"`
	Content   *[]base.ContentItem `json:"content,omitempty"`
	base.Locatable
}

func (e *Composition) Validate() bool {
	//TODO
	return true
}
