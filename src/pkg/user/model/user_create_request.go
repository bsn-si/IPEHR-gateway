package model

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/roles"
)

// Fields Name, Address, Description are required for Doctor role
type UserCreateRequest struct {
	UserID      string `json:"userID"`
	Password    string `json:"password"`
	Role        uint8  `json:"role"`
	Name        string `json:"name,omitempty"`
	Address     string `json:"address,omitempty"`
	Description string `json:"description,omitempty"`
	PictuteURL  string `json:"pictureURL,omitempty"`
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

	if u.Role == uint8(roles.Doctor) {
		switch {
		case u.Name == "":
			return false, errors.ErrFieldIsEmpty("name")
		case u.Address == "":
			return false, errors.ErrFieldIsEmpty("address")
		case u.Description == "":
			return false, errors.ErrFieldIsEmpty("description")
		}
	}

	return true, nil
}
