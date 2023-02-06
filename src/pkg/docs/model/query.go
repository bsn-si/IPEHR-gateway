package model

import "github.com/bsn-si/IPEHR-gateway/src/pkg/aqlprocessor"

// https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/query.html#requirements
type QueryRequest struct {
	Query           string                 `json:"q"`
	Offset          int                    `json:"offset"`
	Fetch           int                    `json:"fetch"`
	QueryParameters map[string]interface{} `json:"query_parameters"`
}

func (q *QueryRequest) Validate() bool {
	if len(q.Query) == 0 {
		return false
	}

	_, err := aqlprocessor.NewAqlProcessor(q.Query).Process()

	return err == nil
}

// https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/query.html#requirements-response-structure
type QueryResponse struct {
	Meta struct {
		Href          string `json:"_href"`
		Type          string `json:"_type"`
		SchemaVersion string `json:"_schema_version"`
		Created       string `json:"_created"`
		Generator     string `json:"_generator"`
		ExecutedAql   string `json:"_executed_aql"`
	} `json:"meta"`
	Name    string        `json:"name"`
	Query   string        `json:"q"`
	Columns []QueryColumn `json:"columns"`
	Rows    []interface{} `json:"rows"`
}

func (q *QueryResponse) Validate() bool {
	//TODO
	return true
}

type QueryColumn struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
