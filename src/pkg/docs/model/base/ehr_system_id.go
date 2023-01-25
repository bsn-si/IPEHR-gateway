package base

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

// EhrSystemID used in a different places like the identifier of the system, e.g. in OBJECT_VERSION_ID this identifier
// defines as a creating_system_id. Typically, it is represented a reverse domain identifier
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_identifying_versions_within_openehr_versioned_containers
type EhrSystemID string

func NewEhrSystemID(ehrSystemID string) (EhrSystemID, error) {
	e := EhrSystemID(ehrSystemID)

	if !e.Validate(ehrSystemID) {
		return "", fmt.Errorf("%w Incorrect system ID", errors.ErrIncorrectFormat)
	}

	return e, nil
}

func (EhrSystemID) Validate(ehrSystemID string) bool {
	validation := true
	if len(ehrSystemID) == 0 {
		validation = false
	}

	return validation
}

func (e EhrSystemID) String() string {
	return string(e)
}

func (e EhrSystemID) Equal(val string) bool {
	return string(e) == val
}
