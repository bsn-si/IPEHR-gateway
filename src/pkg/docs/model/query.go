package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/aqlprocessor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func init() {
	gob.Register(aqlprocessor.IdentifiedPathSelectValue{})
	gob.Register(&big.Int{})
	gob.Register(&big.Float{})
}

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

func (q QueryRequest) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(q)
	if err != nil {
		return nil, fmt.Errorf("QueryRequest gob encode error: %w", err)
	}

	return buf.Bytes(), nil
}

func (q *QueryRequest) FromBytes(data []byte) error {
	var buf = bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(q)
	if err != nil {
		return fmt.Errorf("QueryRequest gob decode error: %w", err)
	}

	return nil
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
