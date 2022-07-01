package model

type SwagComposition struct {
	Type            string `json:"_type"`
	ArchetypeNodeID string `json:"archetype_node_id"`
	Name            struct {
		Value string `json:"value"`
	} `json:"name"`
	UID struct {
		Type  string `json:"_type"`
		Value string `json:"value"`
	} `json:"uid,omitempty"`
	ArchetypeDetails struct {
		Type        string `json:"_type"`
		ArchetypeID struct {
			Value string `json:"value"`
		} `json:"archetype_id"`
		TemplateID struct {
			Value string `json:"value"`
		} `json:"template_id,omitempty"`
		RmVersion string `json:"rm_version"`
	} `json:"archetype_details,omitempty"`
	Language struct {
		TerminologyID struct {
			Value string `json:"value"`
		} `json:"terminology_id"`
		CodeString    string `json:"code_string,omitempty"`
		PreferredTerm string `json:"preferred_term,omitempty"`
	} `json:"language"`
	Territory struct {
		TerminologyID struct {
			Value string `json:"value"`
		} `json:"terminology_id"`
		CodeString    string `json:"code_string,omitempty"`
		PreferredTerm string `json:"preferred_term,omitempty"`
	} `json:"territory"`
	Category struct {
		Value        string `json:"value"`
		DefiningCode struct {
			TerminologyID struct {
				Value string `json:"value"`
			} `json:"terminology_id"`
			CodeString    string `json:"code_string,omitempty"`
			PreferredTerm string `json:"preferred_term,omitempty"`
		} `json:"defining_code"`
	} `json:"category"`
	Composer struct {
		Type        string `json:"_type"`
		Name        string `json:"name"`
		ExternalRef struct {
			ID struct {
				Type  string `json:"_type"`
				Value string `json:"value"`
			} `json:"id"`
			Namespace string `json:"namespace"`
			Type      string `json:"type"`
		} `json:"external_ref"`
	} `json:"composer"`
	Context struct {
		StartTime struct {
			Value string `json:"value"`
		} `json:"start_time"`
		Setting struct {
			Value        string `json:"value"`
			DefiningCode struct {
				TerminologyID struct {
					Value string `json:"value"`
				} `json:"terminology_id"`
				CodeString    string `json:"code_string,omitempty"`
				PreferredTerm string `json:"preferred_term,omitempty"`
			} `json:"defining_code"`
		} `json:"setting"`
	} `json:"context,omitempty"`
	Content []interface{} `json:"content,omitempty"`
} //@name Composition
