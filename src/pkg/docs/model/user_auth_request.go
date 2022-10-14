package model

import (
	"hms/gateway/pkg/errors"

	"github.com/google/uuid"
)

type UserAuthRequest struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

func (u *UserAuthRequest) userUUID() (*uuid.UUID, error) {
	userUUID, err := uuid.Parse(u.UserID)
	if err != nil {
		return nil, err
	}

	return &userUUID, nil
}

func (u *UserAuthRequest) Validate() (bool, error) {
	if len(u.UserID) == 0 {
		return false, errors.ErrFieldIsEmpty("UserId")
	}

	if _, err := u.userUUID(); err != nil {
		return false, errors.ErrFieldIsIncorrect("UserId")
	}

	// TODO Check Password (min max other conds)
	if len(u.Password) == 0 {
		return false, errors.ErrFieldIsEmpty("Password")
	}
	// TODO Check Role
	return true, nil
}