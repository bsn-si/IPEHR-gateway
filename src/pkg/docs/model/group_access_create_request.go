package model

type GroupAccessCreateRequest struct {
	Description string `json:"description"`
}

func (e *GroupAccessCreateRequest) Validate() bool {
	//TODO

	return true
}
