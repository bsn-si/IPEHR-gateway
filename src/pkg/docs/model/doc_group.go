package model

import (
	"hms/gateway/pkg/crypto/chachaPoly"

	"github.com/google/uuid"
)

type DocumentGroup struct {
	GroupID     uuid.UUID `json:"groupID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Members     []string  `json:"members"`

	GroupKey     *chachaPoly.Key `json:"-"`
	GroupKeyEncr []byte          `json:"-"`
	MembersEncr  [][]byte        `json:"-"`
}
