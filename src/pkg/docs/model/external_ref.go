package model

type ExternalRef struct {
	ID struct {
		Value  string `json:"value"`
		Type   string `json:"_type,omitempty"`
		Scheme string `json:"scheme,omitempty"`
	} `json:"id"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
}
