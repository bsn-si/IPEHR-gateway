package model

import "hms/gateway/pkg/docs/model/base"

type EhrSummary struct {
	SystemId  base.ObjectId `json:"system_id"`
	EhrId     base.ObjectId `json:"ehr_id"`
	EhrStatus struct {
		Id        base.ObjectId `json:"id"`
		Namespace string        `json:"namespace"`
		Type      string        `json:"type"`
	} `json:"ehr_status"`
	EhrAccess struct {
		Id        base.ObjectId `json:"id"`
		Namespace string        `json:"namespace"`
		Type      string        `json:"type"`
	} `json:"ehr_access"`
	TimeCreated struct {
		Value string `json:"value"`
	} `json:"time_created"`
}
