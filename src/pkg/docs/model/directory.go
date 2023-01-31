package model

import (
	"strings"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type Directory struct {
	base.Locatable
	FeederAudit base.FeederAudit   `json:"feeder_audit"`
	Folders     []*Directory       `json:"folders"`
	Items       []DirectoryItem    `json:"items,omitempty"`
	Details     base.ItemStructure `json:"details"`
}

type DirectoryItem struct {
	ID        base.UIDBasedID `json:"id"`
	Type      base.ItemType   `json:"type"`
	Namespace string          `json:"namespace"`
}

const (
	directorySeparator      = "/"
	directoryRootFolderName = "root"
)

func (d *Directory) Validate() error {
	var err error

	var errs []error

	if d.UID == nil || d.UID.Value == "" {
		errs = append(errs, errors.ErrFieldIsEmpty("uid"))
	}

	// TODO add ITEMS validation

	for i, e := range errs {
		if i == 0 {
			err = e
			continue
		}

		err = errors.Wrap(err, e.Error())
	}

	return err
}

func (d *Directory) GetByPath(p string) (*Directory, error) {
	p = d.sanitize(p)
	if p == "" {
		p = directoryRootFolderName
	}

	paths := strings.SplitN(p, directorySeparator, 2)

	if d.Name.Value != paths[0] {
		return nil, errors.ErrNotFound
	}

	if len(paths) == 2 {
		for _, dd := range d.Folders {
			if dt, err := dd.GetByPath(paths[1]); err == nil {
				return dt, nil
			}
		}

		return nil, errors.ErrNotFound
	}

	return d, nil
}

func (d *Directory) sanitize(p string) string {
	return strings.Trim(p, directorySeparator)
}
