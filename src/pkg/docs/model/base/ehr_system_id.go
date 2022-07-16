package base

import (
	"fmt"
	"hms/gateway/pkg/errors"
)

// EhrSystemID used in a different places like the identifier of the system, e.g. in OBJECT_VERSION_ID this identifier
// defines as a creating_system_id. Typically, it is represented a reverse domain identifier
// https://specifications.openehr.org/releases/BASE/latest/base_types.html#_identifying_versions_within_openehr_versioned_containers
type EhrSystemID struct {
	ehrSystemID string
}

func NewEhrSystemID(ehrSystemID string) (*EhrSystemID, error) {
	e := &EhrSystemID{
		ehrSystemID: ehrSystemID,
	}

	if !e.Validate(ehrSystemID) {
		return nil, fmt.Errorf("%w Incorrect system ID", errors.ErrIncorrectFormat)
	}

	e.ehrSystemID = ehrSystemID

	return e, nil
}

func (*EhrSystemID) Validate(ehrSystemID string) bool {
	validation := true
	if len(ehrSystemID) == 0 {
		validation = false
	}

	return validation
}

func (o *EhrSystemID) String() string {
	return o.ehrSystemID
}

func (o *EhrSystemID) Equal(val string) bool {
	return val == o.ehrSystemID
}
