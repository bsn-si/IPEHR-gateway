package model

import "hms/gateway/pkg/docs/model/base"

// AQL stored query
// https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/StoredQuery_schema

type StoredQuery struct {
	Name        base.QueryName  `json:"name"`
	Type        base.QueryType  `json:"type"`
	Version     string          `json:"version,omitempty"`
	TimeCreated base.DvDateTime `json:"saved,omitempty"`
	Query       string          `json:"q"`
}

func (q *StoredQuery) Validate() bool {
	//TODO
	return true
}
