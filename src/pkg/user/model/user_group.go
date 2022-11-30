package model

import "github.com/google/uuid"

type UserGroup struct {
	GroupID     *uuid.UUID `json:"-"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}
