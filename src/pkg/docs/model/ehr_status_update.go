package model

type EhrStatusUpdate struct {
	Type            string `json:"_type"`
	ArchetypeNodeID string `json:"archetype_node_id"`
	Name            struct {
		Value string `json:"value"`
	} `json:"name"`
	UID struct {
		Type  string `json:"_type"`
		Value string `json:"value"`
	} `json:"uid"`
	Subject struct {
		ExternalRef ExternalRef `json:"external_ref"`
	} `json:"subject"`
	OtherDetails struct {
		Type            string `json:"_type"`
		ArchetypeNodeID string `json:"archetype_node_id"`
		Name            struct {
			Value string `json:"value"`
		} `json:"name"`
		Items []interface{} `json:"items"`
	} `json:"other_details"`
	IsModifiable bool `json:"is_modifiable"`
	IsQueryable  bool `json:"is_queryable"`
}
