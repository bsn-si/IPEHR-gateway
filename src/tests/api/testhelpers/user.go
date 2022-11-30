package testhelpers

import "hms/gateway/pkg/user/model"

type UserHelper struct {
}

type AuthOption func(*model.UserAuthRequest)

func (UserHelper) UserAuthRequest(options ...AuthOption) *model.UserAuthRequest {
	p := &model.UserAuthRequest{
		UserID:   "",
		Password: "",
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func (UserHelper) WithPassword(val string) AuthOption {
	return func(r *model.UserAuthRequest) {
		r.Password = val
	}
}

func (UserHelper) WithUserID(val string) AuthOption {
	return func(r *model.UserAuthRequest) {
		r.UserID = val
	}
}
