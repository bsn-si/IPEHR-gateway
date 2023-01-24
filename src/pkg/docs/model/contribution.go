package model

import (
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	errorsPkg "github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/helper"
)

// AUDIT_DETAILS
// The set of attributes required to document the committal of an information item to a repository.
// https://specifications.openehr.org/releases/RM/latest/common.html#_audit_details_class
type AuditDetails struct {
	Type          base.ItemType    `json:"_type"`
	SystemID      string           `json:"system_id"`
	TimeCommitted base.DvDateTime  `json:"time_committed,omitempty"`
	ChangeType    base.DvCodedText `json:"change_type,omitempty"`
	Committer     base.PartyProxy  `json:"committer,omitempty"`
	Description   base.DvText      `json:"description,omitempty"`
}

// https://specifications.openehr.org/releases/RM/latest/common.html#_contribution_class
type Contribution struct {
	base.Root
	UID      base.UIDBasedID       `json:"uid"`
	Versions []ContributionVersion `json:"versions"`
	Audit    AuditDetails          `json:"audit"`
}

type ContributionResponse struct {
	UID      base.UIDBasedID               `json:"uid"`
	Versions []ContributionVersionResponse `json:"versions"`
	Audit    AuditDetails                  `json:"audit"`
}

type ContributionVersion struct {
	Type                base.ItemType    `json:"_type"`
	Contribution        base.ObjectRef   `json:"contribution"`
	CommitAudit         AuditDetails     `json:"commit_audit"`
	UID                 base.UIDBasedID  `json:"uid"`
	PrecedingVersionUID base.UIDBasedID  `json:"preceding_version_uid"`
	LifecycleState      base.DvCodedText `json:"lifecycle_state"`
	Data                base.Root        `json:"data"`
}

type ContributionVersionResponse struct {
	Type         base.ItemType `json:"_type"`
	Contribution base.ObjectRef
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
	case base.FolderItemType:
		d := Directory{}
		if err := json.Unmarshal(data, &d); err != nil {
			return errors.Wrapf(err, "cannot unmarshal contribution content item: '%v'", tmp.Type)
		}

		w.item = d
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

// TODO if type contribution version was modified then check that version is exist
func (c *Contribution) Validate(template helper.Searcher) (bool, error) {
	if len(c.Versions) == 0 {
		return false, errorsPkg.ErrFieldIsEmpty("Versions")
	}

	for _, v := range c.Versions {
		if ok, err := v.Validate(template); !ok {
			return false, errorsPkg.Wrap(err, "Version is invalid")
		}
	}

	return true, nil
}

func (c *ContributionVersion) Validate(templateSearcher helper.Searcher) (bool, error) {
	if c.Data == nil {
		return false, errorsPkg.ErrFieldIsEmpty("Data")
	}

	allowedVersions := []base.ItemType{base.VersionOriginalItemType, base.VersionImportedItemType}
	if !slices.Contains(allowedVersions, c.Type) {
		return false, errorsPkg.ErrTypeNotValid
	}

	switch c.Data.GetType() {
	case base.CompositionItemType:
		composition := c.Data.(Composition)

		// If version of lifecycle state is incomplete then validation can be missed partially
		// https://specifications.openehr.org/releases/RM/latest/common.html#_version_lifecycle
		if ok, err := composition.Validate(); !ok {
			return false, errorsPkg.Wrap(err, "Version of contribution is not valid")
		}

		templateID := composition.ArchetypeDetails.TemplateID.Value
		if templateID == "" {
			return false, errorsPkg.ErrFieldIsEmpty("TemplateID")
		}

		ok, err := templateSearcher.IsExist(templateID)
		if !ok || err != nil {
			return false, errorsPkg.ErrObjectWithIDIsNotExist("TemplateID", templateID)
		}
	}

	return true, nil
}
