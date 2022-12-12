package model

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model/base"

	"github.com/pkg/errors"
)

// https://specifications.openehr.org/releases/RM/latest/common.html#_contribution_class
type Contribution struct {
	UID base.UIDBasedID `json:"uid"`
	//Versions *[]base.ObjectRef `json:"versions"`
	Versions []ContributionVersion `json:"versions"`
	Audit    AuditDetails          `json:"audit"`
}

//type ContributionResponse struct{}

type ContributionVersion struct {
	Type           base.ItemType    `json:"_type"`
	Contribution   base.ObjectRef   `json:"contribution"`
	CommitAudit    AuditDetails     `json:"commit_audit"`
	UID            base.UIDBasedID  `json:"uid"`
	Data           base.Root        `json:"data"`
	LifecycleState base.DvCodedText `json:"lifecycle_state"`
	//Data           interface{}      `json:"data"` // TODO or *base.PartyProxy ???, or generics
	//Data           *MainType[T]     `json:"data"` // TODO or *base.PartyProxy ???, or generics
	//Data *base.dataValueWrapper `json:"data"`
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
	Type         base.ItemType    `json:"_type"`
	SystemID     string           `json:"system_id"`
	TimeCommited base.DvDateTime  `json:"time_committed,omitempty"`
	ChangeType   base.DvCodedText `json:"change_type,omitempty"`
	Committer    base.PartyProxy  `json:"committer,omitempty"`
	Description  base.DvText      `json:"description,omitempty"`
}

func (c *Contribution) UnmarshalJSON(data []byte) error {
	cc := contributionWrapper{}
	if err := json.Unmarshal(data, &cc); err != nil {
		return errors.Wrap(err, "cannot unmarshal 'contribution' struct from json bytes")
	}

	//c.Type = cc.Type
	//c.Name = cc.Name

	//if cc.Content != nil {
	//	c.Content = make([]base.Root, 0, len(cc.Content))
	//	for _, item := range cc.Content {
	//		c.Content = append(c.Content, item.item)
	//	}
	//}

	return nil
}

type contributionWrapper struct {
	Type           base.ItemType              `json:"_type"`
	Contribution   base.ObjectRef             `json:"contribution"`
	CommitAudit    AuditDetails               `json:"commit_audit"`
	UID            base.UIDBasedID            `json:"uid"`
	LifecycleState base.DvCodedText           `json:"lifecycle_state"`
	Data           contributionContentWrapper `json:"data"`
}

type contributionContentWrapper struct {
	item base.Root
}

func (w *contributionContentWrapper) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type base.ItemType `json:"_type"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "can't unmarshal contribution content wrapper")
	}

	switch tmp.Type {
	case base.CompositionItemType:
		w.item = &Composition{}
	default:
		return errors.Errorf("unexpected contribution content item: '%v'", tmp.Type)
	}

	if err := json.Unmarshal(data, w.item); err != nil {
		return errors.Wrapf(err, "cannot unmarshal contribution content item: '%v'", tmp.Type)
	}

	return nil
}
