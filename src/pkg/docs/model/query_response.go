package model

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
	Name    string `json:"name"`
	Query   string `json:"q"`
	Columns []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"columns"`
	Rows []interface{} `json:"rows"`
}

func (q *QueryResponse) Validate() bool {
	//TODO
	return true
}
