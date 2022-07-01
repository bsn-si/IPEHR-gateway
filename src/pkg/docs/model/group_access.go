package model

import (
	"github.com/google/uuid"

	"hms/gateway/pkg/crypto/chachaPoly"
)

type GroupAccess struct {
	GroupUUID   *uuid.UUID      `json:"group_id"`
	Description string          `json:"description"`
	Key         *chachaPoly.Key `json:"-"`
	Nonce       *[12]byte       `json:"-"`
}
