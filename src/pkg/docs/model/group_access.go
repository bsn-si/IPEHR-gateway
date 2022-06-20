package model

import "hms/gateway/pkg/crypto/chacha_poly"

type GroupAccess struct {
	GroupId     string           `json:"group_id"`
	Description string           `json:"description"`
	Key         *chacha_poly.Key `json:"-"`
}
