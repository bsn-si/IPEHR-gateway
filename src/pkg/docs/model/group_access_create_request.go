package model

type GroupAccessCreateRequest struct {
	Description string `json:"description"`
}

// TODO
func (e *GroupAccessCreateRequest) Validate() bool {
	return true
}
