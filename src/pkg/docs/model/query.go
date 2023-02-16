package model

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aqlprocessor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

// https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/query.html#requirements
type QueryRequest struct {
	Query           string                 `json:"q"`
	Offset          int                    `json:"offset"`
	Fetch           int                    `json:"fetch"`
	QueryParameters map[string]interface{} `json:"query_parameters"`
	QueryParsed     *aqlprocessor.Query    `json:"-"`
}

func (q *QueryRequest) Validate() error {
	if len(q.Query) == 0 {
		return errors.ErrFieldIsEmpty("query")
	}

	return nil
}

func (q *QueryRequest) AqlProcess() error {
	var err error

	if q.Offset != 0 {
		q.Query = fmt.Sprintf("%s OFFSET %d", q.Query, q.Offset)
	}

	if q.Fetch != 0 {
		q.Query = fmt.Sprintf("%s LIMIT %d", q.Query, q.Fetch)
	}

	q.QueryParsed, err = aqlprocessor.NewAqlProcessor(q.Query).Process()
	return err
}

func (q *QueryRequest) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal(q.QueryParsed)
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
