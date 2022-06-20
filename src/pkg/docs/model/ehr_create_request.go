package model

type EhrCreateRequest struct {
	Type            string `json:"_type"`
	ArchetypeNodeId string `json:"archetype_node_id"`
	Name            struct {
		Value string `json:"value"`
	} `json:"name"`
	Subject struct {
		ExternalRef ExternalRef `json:"external_ref"`
	} `json:"subject"`
	IsModifiable bool `json:"isModifiable"`
	IsQueryable  bool `json:"isQueryable"`
}

func (e *EhrCreateRequest) Validate() bool {
	//TODO

	return true
}
