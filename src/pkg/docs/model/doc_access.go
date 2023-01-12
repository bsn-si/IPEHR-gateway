package model

type DocAccessDocument struct {
	CID         string `json:"CID"`
	Level       string `json:"level"`
	Description string `json:"description"`
}

type DocAccessDocumentGroup struct {
	GroupID       string   `json:"groupID"`
	Level         string   `json:"level"`
	Description   string   `json:"description"`
	ParentGroupID string   `json:"parentGroupID,omitempty"`
	Documents     []string `json:"documents"`
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
