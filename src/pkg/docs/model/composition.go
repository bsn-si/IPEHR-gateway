package model

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model/base"

	"github.com/pkg/errors"
)

// Composition Content of one version in a VERSIONED_COMPOSITION. A Composition is considered the unit
// of modification of the record, the unit of transmission in record Extracts, and the unit of
// attestation by authorising clinicians. In this latter sense, it may be considered equivalent to a
// signed document.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_composition_class
type Composition struct {
	Language  base.CodePhrase  `json:"language"`
	Territory base.CodePhrase  `json:"territory"`
	Category  base.DvCodedText `json:"category"`
	Context   *EventContext    `json:"context,omitempty"`
	Composer  base.PartyProxy  `json:"composer"`
	Content   []base.Section   `json:"content,omitempty"`
	base.Locatable
}

func (c *Composition) Validate() bool {
	validation := true
	if c.Type != base.CompositionItemType {
		validation = false
	}

	return validation
}

type composition struct {
	Language  base.CodePhrase  `json:"language"`
	Territory base.CodePhrase  `json:"territory"`
	Category  base.DvCodedText `json:"category"`
	Context   *EventContext    `json:"context,omitempty"`
	Composer  base.PartyProxy  `json:"composer"`
	Content   []base.Section   `json:"content,omitempty"`
	base.Locatable
}

func (c *Composition) UnmarshalJSON(data []byte) error {
	cc := composition{}
	if err := json.Unmarshal(data, &cc); err != nil {
		return errors.Wrap(err, "cannot unmarshal 'composition' struct from json bytes")
	}

	c.Type = cc.Type
	c.Name = cc.Name
	c.ArchetypeNodeID = cc.ArchetypeNodeID
	c.ObjectVersionID = cc.ObjectVersionID
	c.Composer = cc.Composer
	c.Context = cc.Context
	c.Category = cc.Category
	c.Territory = cc.Territory
	c.Language = cc.Language
	c.Pathable = cc.Pathable
	c.ArchetypeDetails = cc.ArchetypeDetails
	c.Links = cc.Links
	c.Content = cc.Content

	return nil
}
