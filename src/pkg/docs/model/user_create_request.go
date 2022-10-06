package model

type UserCreateRequest struct {
	SystemID string `json:"systemID"`
	UserID   string `json:"userID"`
	Password string `json:"password"`
	Role     uint8  `json:"role"`
}

func (e *UserCreateRequest) Validate() bool {
	// TODO Check userId
	// TODO Check Password
	// TODO Check SystemId
	// TODO Check Role
	return true
}
