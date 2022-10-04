package model

type UserCreateRequest struct {
	SystemID string `json:"systemID"`
	UserID   string `json:"userID"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (e *UserCreateRequest) Validate() bool {
	// TODO Check userId
	// TODO Check Password
	// TODO Check SystemId
	// TODO Check Role
	return true
}
