package model

import (
	"hms/gateway/pkg/errors"
)

type UserCreateRequest struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
	Role     uint8  `json:"role"`
}

func (u *UserCreateRequest) Validate() (bool, error) {
	if len(u.UserID) == 0 {
		return false, errors.ErrFieldIsEmpty("UserId")
	}

	//TODO format validate

	// TODO Check Password (min max other conds)
	if len(u.Password) == 0 {
		return false, errors.ErrFieldIsEmpty("Password")
	}
	// TODO Check Role
	return true, nil
}
