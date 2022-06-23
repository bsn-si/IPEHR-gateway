package parser

import (
	"encoding/json"

	"hms/gateway/pkg/docs/model"
)

func ParseEhr(inDocument []byte) (doc *model.EHR, err error) {
	doc = &model.EHR{}
	return json.Unmarshal(inDocument, doc)
}
