package model

import (
	"github.com/google/uuid"

	"hms/gateway/pkg/crypto/chacha_poly"
)

type GroupAccess struct {
	GroupUUID   *uuid.UUID       `json:"group_id"`
	Description string           `json:"description"`
	Key         *chacha_poly.Key `json:"-"`
}
