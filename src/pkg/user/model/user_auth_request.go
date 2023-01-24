package model

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type UserAuthRequest struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

func (u *UserAuthRequest) Validate() (bool, error) {
	if len(u.UserID) == 0 {
		return false, errors.ErrFieldIsEmpty("UserId")
	}

	// TODO Check Password (min max other conds)
	if len(u.Password) == 0 {
		return false, errors.ErrFieldIsEmpty("Password")
	}
	// TODO Check Role
	return true, nil
}
