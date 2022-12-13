package model

type EhrDocumentItem struct {
	Name        string `json:"name"`
	UID         string `json:"uid"`
	TimeCreated string `json:"timeCreated"`
}
