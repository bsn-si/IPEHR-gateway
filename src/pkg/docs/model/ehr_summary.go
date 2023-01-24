package model

import "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"

type EhrSummary struct {
	SystemID  base.ObjectID `json:"system_id"`
	EhrID     base.ObjectID `json:"ehr_id"`
	EhrStatus struct {
		ID        base.ObjectID `json:"id"`
		Namespace string        `json:"namespace"`
		Type      string        `json:"type"`
	} `json:"ehr_status"`
	EhrAccess struct {
		ID        base.ObjectID `json:"id"`
		Namespace string        `json:"namespace"`
		Type      string        `json:"type"`
	} `json:"ehr_access"`
	TimeCreated struct {
		Value string `json:"value"`
	} `json:"time_created"`
}
