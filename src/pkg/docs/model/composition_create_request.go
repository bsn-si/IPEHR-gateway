package model

import "time"

type CompositionCreateRequest struct {
	// TODO поставить реальные
	Type            string `json:"_type"`
	ArchetypeNodeID string `json:"archetype_node_id"`
	Name            struct {
		Value string `json:"value"`
	} `json:"name"`
	UID struct {
		Type  string `json:"_type"`
		Value string `json:"value"`
	} `json:"uid"`
	ArchetypeDetails struct {
		ArchetypeID struct {
			Value string `json:"value"`
		} `json:"archetype_id"`
		TemplateID struct {
			Value string `json:"value"`
		} `json:"template_id"`
		RmVersion string `json:"rm_version"`
	} `json:"archetype_details"`
	Language struct {
		TerminologyID struct {
			Value string `json:"value"`
		} `json:"terminology_id"`
		CodeString string `json:"code_string"`
	} `json:"language"`
	Territory struct {
		TerminologyID struct {
			Value string `json:"value"`
		} `json:"terminology_id"`
		CodeString string `json:"code_string"`
	} `json:"territory"`
	Category struct {
		Value        string `json:"value"`
		DefiningCode struct {
			TerminologyID struct {
				Value string `json:"value"`
			} `json:"terminology_id"`
			CodeString string `json:"code_string"`
		} `json:"defining_code"`
	} `json:"category"`
	Composer struct {
		Type        string `json:"_type"`
		ExternalRef struct {
			ID struct {
				Type  string `json:"_type"`
				Value string `json:"value"`
			} `json:"id"`
			Namespace string `json:"namespace"`
			Type      string `json:"type"`
		} `json:"external_ref"`
		Name string `json:"name"`
	} `json:"composer"`
	Context struct {
		StartTime struct {
			Value time.Time `json:"value"`
		} `json:"start_time"`
		Setting struct {
			Value        string `json:"value"`
			DefiningCode struct {
				TerminologyID struct {
					Value string `json:"value"`
				} `json:"terminology_id"`
				CodeString string `json:"code_string"`
			} `json:"defining_code"`
		} `json:"setting"`
	} `json:"context"`
	Content []interface{} `json:"content"`
}

func (e *CompositionCreateRequest) Validate() bool {
	//TODO
	return true
}
