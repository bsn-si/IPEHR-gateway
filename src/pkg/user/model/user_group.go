package model

import "github.com/google/uuid"

type UserGroup struct {
	GroupID     *uuid.UUID `json:"groupID"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Members     []string   `json:"members"`

	GroupKey     [32]byte `json:"-"`
	GroupKeyEncr []byte   `json:"-"`
	ContentEncr  []byte   `json:"-"`
	MembersEncr  [][]byte `json:"-"`
}
