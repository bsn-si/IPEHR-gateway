package model

import "hms/gateway/pkg/docs/model/base"

type ExternalRef struct {
	Id        base.ObjectId `json:"id"`
	Namespace string        `json:"namespace"`
	Type      string        `json:"type"`
	Scheme    string        `json:"scheme,omitempty"`
}
