package base

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Section
// Represents a heading in a heading structure, or section tree.
// Created according to archetyped structures for typical headings such as SOAP, physical examination,
// but also pathology result heading structures. Should not be used instead of ENTRY hierarchical structures.
// https://specifications.openehr.org/releases/RM/Release-1.0.2/ehr.html#_section_class
type Section struct {
	Locatable
	Items []Root `json:"items,omitempty"`
}

func (s *Section) UnmarshalJSON(data []byte) error {
	ss := sectionWrapper{}
	if err := json.Unmarshal(data, &ss); err != nil {
		return errors.Wrap(err, "cannot unmarshal section")
	}

	s.Locatable = ss.Locatable
	s.Items = make([]Root, 0, len(ss.Items))

	for _, item := range ss.Items {
		s.Items = append(s.Items, item.contentItem)
	}

	return nil
}

type sectionWrapper struct {
	Locatable
	Items []sectionItemWrapper `json:"items,omitempty"`
}

type sectionItemWrapper struct {
	contentItem Root `json:"-"`
}

func (item *sectionItemWrapper) UnmarshalJSON(data []byte) error {
	str := struct {
		Type ItemType `json:"_type"`
	}{}
	if err := json.Unmarshal(data, &str); err != nil {
		return errors.Wrap(err, "cannot unmarshal section item wrapper")
	}

	switch str.Type {
	case ActionItemType:
		item.contentItem = &Action{}
	case EvaluationItemType:
		item.contentItem = &Evaluation{}
	case ObservationItemType:
		item.contentItem = &Observation{}
	case InstructionItemType:
		item.contentItem = &Instruction{}
	default:
		return errors.Errorf("unexpected section item type: '%v'", str.Type)
	}

	if err := json.Unmarshal(data, item.contentItem); err != nil {
		return errors.Wrapf(err, "cannot unmarshal secion item type: '%v'", str.Type)
	}

	return nil
}
