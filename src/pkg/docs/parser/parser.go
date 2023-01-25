package parser

import (
	"encoding/json"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
)

func ParseEhr(inDocument []byte) (doc *model.EHR, err error) {
	doc = &model.EHR{}
	err = json.Unmarshal(inDocument, doc)

	return
}
