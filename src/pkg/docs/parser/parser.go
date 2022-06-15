package parser

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model"
)

func ParseDocument(inDocument []byte) (doc model.EHR, err error) {
	err = json.Unmarshal(inDocument, &doc)
	return
}

func ParseComposition(inComposition []byte) (composition model.Composition, err error) {
	err = json.Unmarshal(inComposition, &composition)
	return
}
