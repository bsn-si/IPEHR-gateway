package model

import (
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
)

type Directory struct {
	base.Locatable
	FeederAudit base.FeederAudit   `json:"feeder_audit"`
	Folders     []Directory        `json:"folders"`
	Items       []DirectoryItem    `json:"items,omitempty"`
	Details     base.ItemStructure `json:"details"`
}

type DirectoryItem struct {
	ID        base.UIDBasedID `json:"id"`
	Type      base.ItemType   `json:"type"`
	Namespace string          `json:"namespace"`
}

// TODO
func (d *Directory) GetByPath(p string) (*Directory, error) {
	return nil, errors.ErrNotImplemented
}
