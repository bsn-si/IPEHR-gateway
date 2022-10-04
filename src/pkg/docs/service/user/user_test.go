package user_test

import (
	"hms/gateway/pkg/docs/service/user"
	"testing"
)

func TestService_Register(t *testing.T) {
	service := new(user.Service)
	_ = service.Register("a", "b", "c", "d")
}
