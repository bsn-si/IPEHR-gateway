package model

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

// AQL stored query
// https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/StoredQuery_schema

type StoredQuery struct {
	Name        QueryName `json:"name"`
	Type        QueryType `json:"type"`
	Version     string    `json:"version"`
	TimeCreated string    `json:"saved"`
	Query       string    `json:"q"`
}

func (q *StoredQuery) Validate() error {
	var err error

	var errs []error

	if q.Name == "" {
		errs = append(errs, errors.ErrFieldIsEmpty("name"))
	}

	if q.Type == "" {
		errs = append(errs, errors.ErrFieldIsEmpty("type"))
	}

	if q.Query == "" {
		errs = append(errs, errors.ErrFieldIsEmpty("query"))
	}

	if q.Version == "" {
		errs = append(errs, errors.ErrFieldIsEmpty("version"))
	}

	if q.TimeCreated == "" {
		errs = append(errs, errors.ErrFieldIsEmpty("timeCreated"))
	}

	for i, e := range errs {
		if i == 0 {
			err = e
			continue
		}

		err = errors.Wrap(err, e.Error())
	}

	return err
}

// Query formalism type
type QueryType = string

const (
	QueryTypeAQL QueryType = "AQL"
)

// QueryName
// The (fully qualified) name of the query (when is registered as a stored query), in a format of [{namespace}::]{query-name}. The namespace prefix is optional, and when used it should be in a form of a reverse domain name.
// https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/StoredQuery_schema
type QueryName = string
