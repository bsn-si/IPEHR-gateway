package model

type UserInfo struct {
	Role        string `json:"role"`
	Name        string `json:"name,omitempty"`
	Address     string `json:"address,omitempty"`
	Description string `json:"description,omitempty"`
	PictuteURL  string `json:"pictureURL,omitempty"`
	Code        string `json:"code,omitempty"`
	TimeCreated string `json:"timeCreated"`

	Timestamp []byte `json:"-" msgpack:"timestamp"`
}
