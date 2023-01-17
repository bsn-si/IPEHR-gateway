package model

type DocAccessDocument struct {
	CID         string `json:"CID"`
	Level       string `json:"level,omitempty"`
	Description string `json:"description,omitempty"`
}

type DocAccessDocumentGroup struct {
	GroupID       string               `json:"groupID"`
	Level         string               `json:"level"`
	Description   string               `json:"description"`
	ParentGroupID string               `json:"parentGroupID,omitempty"`
	Documents     []*DocAccessDocument `json:"documents"`
}

type DocAccessSetRequest struct {
	UserID      string
	CID         string
	AccessLevel string
}

type DocAccessListResponse struct {
	Documents      []*DocAccessDocument      `json:"documents"`
	DocumentGroups []*DocAccessDocumentGroup `json:"documentGroups"`
}
