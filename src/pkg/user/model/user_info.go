package model

import "github.com/google/uuid"

type UserInfo struct {
	Role        string     `json:"role"`
	UserID      string     `json:"userID,omitempty"`
	Name        string     `json:"name,omitempty"`
	Address     string     `json:"address,omitempty"`
	Description string     `json:"description,omitempty"`
	PictuteURL  string     `json:"pictureURL,omitempty"`
	Code        string     `json:"code,omitempty"`
	TimeCreated string     `json:"timeCreated"`
	EhrID       *uuid.UUID `json:"ehrID,omitempty"`

	Timestamp []byte `json:"-" msgpack:"timestamp"`
}
