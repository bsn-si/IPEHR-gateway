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
	UID      base.UIDBasedID         `json:"uid"`
	Versions []ContributionVersioner `json:"versions"`
	Audit    AuditDetails            `json:"audit"`
}

type contributionWrapper struct {
	UID      base.UIDBasedID              `json:"uid"`
	Versions []contributionVersionWrapper `json:"versions"`
	Audit    AuditDetails                 `json:"audit"`
}

type ContributionVersioner interface {
	GetVersion() ContributionVersion
}

type ContributionVersion struct {
	Type           base.ItemType    `json:"_type"`
	Contribution   base.ObjectRef   `json:"contribution"`
	CommitAudit    AuditDetails     `json:"commit_audit"`
	UID            base.UIDBasedID  `json:"uid"`
	Data           base.Root        `json:"data"`
	LifecycleState base.DvCodedText `json:"lifecycle_state"`
}

func (v ContributionVersion) GetVersion() ContributionVersion {
	return v
}

type contributionVersionWrapper struct {
	ContributionVersion
	Data contributionVersionDataWrapper `json:"data"`
}

type contributionVersionDataWrapper struct {
	item base.Root
}

func (c *Contribution) UnmarshalJSON(data []byte) error {
	cW := contributionWrapper{}
	if err := json.Unmarshal(data, &cW); err != nil {
		return errors.Wrap(err, "cannot unmarshal 'contribution' struct from json bytes")
	}

	c.UID = cW.UID
	c.Audit = cW.Audit

	if cW.Versions != nil {
		c.Versions = make([]ContributionVersioner, 0, len(cW.Versions))
		for _, item := range cW.Versions {
			c.Versions = append(c.Versions, item.GetVersion())
		}
	}

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
		w.item = &Composition{}
	default:
		return errors.Errorf("unexpected contribution content item: '%v'", tmp.Type)
	}

	if err := json.Unmarshal(data, w.item); err != nil {
		return errors.Wrapf(err, "cannot unmarshal contribution content item: '%v'", tmp.Type)
	}

	return nil
}

func (cV *Contribution) Validate() {
	//TODO invoke validation in ContributionVersion by loop
}
func (cV *ContributionVersion) Validate() {
	// TODO data should exist and type is known, and also run validation in there, if it modifyed and type not exist...
}
