package model

import "github.com/google/uuid"

type UserGroup struct {
	GroupID     *uuid.UUID `json:"groupID"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Members     []string   `json:"members"`

	Key         [32]byte `json:"-" msgpack:"-"`
	KeyEncr     []byte   `json:"-" msgpack:"-"`
	IDEncr      []byte   `json:"-" msgpack:"-"`
	ContentEncr []byte   `json:"-" msgpack:"-"`
	MembersEncr [][]byte `json:"-" msgpack:"-"`
	Packed      []byte   `json:"-" msgpack:"-"`
}
