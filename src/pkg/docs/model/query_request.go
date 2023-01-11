package model

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

	return true
}
