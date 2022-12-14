package model

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model/base"

	"github.com/pkg/errors"
)

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

// https://specifications.openehr.org/releases/RM/latest/common.html#_contribution_class
type Contribution struct {
	UID      base.UIDBasedID       `json:"uid"`
	Versions []ContributionVersion `json:"versions"`
	Audit    AuditDetails          `json:"audit"`
}

type ContributionVersion struct {
	Type           base.ItemType    `json:"_type"`
	Contribution   base.ObjectRef   `json:"contribution"`
	CommitAudit    AuditDetails     `json:"commit_audit"`
	UID            base.UIDBasedID  `json:"uid"`
	LifecycleState base.DvCodedText `json:"lifecycle_state"`
	Data           base.Root        `json:"data"`
}

type contributionVersionWrapper struct {
	Type           base.ItemType                  `json:"_type"`
	Contribution   base.ObjectRef                 `json:"contribution"`
	CommitAudit    AuditDetails                   `json:"commit_audit"`
	UID            base.UIDBasedID                `json:"uid"`
	LifecycleState base.DvCodedText               `json:"lifecycle_state"`
	Data           contributionVersionDataWrapper `json:"data"`
}

type contributionVersionDataWrapper struct {
	item base.Root
}

func (c *ContributionVersion) UnmarshalJSON(data []byte) error {
	w := contributionVersionWrapper{}
	if err := json.Unmarshal(data, &w); err != nil {
		return errors.Wrap(err, "cannot unmarshal 'contribution' struct from json bytes")
	}

	c.UID = w.UID
	c.Data = w.Data.item
	c.LifecycleState = w.LifecycleState
	c.CommitAudit = w.CommitAudit
	c.Contribution = w.Contribution
	c.Type = w.Type

	return nil
}

func (w *contributionVersionDataWrapper) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type base.ItemType `json:"_type"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "can't unmarshal contribution content wrapper")
	}

	switch tmp.Type {
	case base.CompositionItemType:
		c := Composition{}
		if err := c.UnmarshalJSON(data); err != nil {
			return errors.Wrapf(err, "cannot unmarshal contribution content item: '%v'", tmp.Type)
		}
		w.item = c
	default:
		return errors.Errorf("unexpected contribution content item: '%v'", tmp.Type)
	}

	return nil
}

func (cV *Contribution) Validate() {
	//TODO invoke validation in ContributionVersion by loop
}
func (cV *ContributionVersion) Validate() {
	// TODO data should exist and type is known, and also run validation in there, if it modifyed and type not exist...
}