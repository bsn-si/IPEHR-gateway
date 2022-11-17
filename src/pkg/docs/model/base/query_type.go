package base

//Query formalism type
//https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/StoredQuery_schema

type QueryType string

func (qt QueryType) ToString() string {
	return string(qt)
}

const (
	AQLQueryType QueryType = "AQL"
)
