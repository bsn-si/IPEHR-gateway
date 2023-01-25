package model

import (
	"encoding/json"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	errorsPkg "github.com/bsn-si/IPEHR-gateway/src/pkg/errors"

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
	Content   []base.Root      `json:"content,omitempty"`
	base.Locatable
}

func (c *Composition) Validate() (bool, error) {
	if c.Type != base.CompositionItemType {
		return false, errorsPkg.ErrIsUnsupported
	}

	return true, nil
}

func (c *Composition) UnmarshalJSON(data []byte) error {
	cc := compositionWrapper{}
	if err := json.Unmarshal(data, &cc); err != nil {
		return errors.Wrap(err, "cannot unmarshal 'composition' struct from json bytes")
	}

	c.Type = cc.Type
	c.Name = cc.Name
	c.ArchetypeDetails = cc.ArchetypeDetails
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

	if cc.Content != nil {
		c.Content = make([]base.Root, 0, len(cc.Content))
		for _, item := range cc.Content {
			c.Content = append(c.Content, item.item)
		}
	}

	return nil
}

type compositionWrapper struct {
	Language  base.CodePhrase             `json:"language"`
	Territory base.CodePhrase             `json:"territory"`
	Category  base.DvCodedText            `json:"category"`
	Context   *EventContext               `json:"context,omitempty"`
	Composer  base.PartyProxy             `json:"composer"`
	Content   []compositionContentWrapper `json:"content,omitempty"`
	base.Locatable
}

type compositionContentWrapper struct {
	item base.Root
}

func (w *compositionContentWrapper) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type base.ItemType `json:"_type"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "can't unmarshal composition content wrapper")
	}

	switch tmp.Type {
	case base.SectionItemType:
		fallthrough
	case base.EvaluationItemType:
		w.item = &base.Section{}
	default:
		return errors.Errorf("unexpected composition content item: '%v'", tmp.Type)
	}

	if err := json.Unmarshal(data, w.item); err != nil {
		return errors.Wrapf(err, "cannot unmarshal composition content item: '%v'", tmp.Type)
	}

	return nil
}
