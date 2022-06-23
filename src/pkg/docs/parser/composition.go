package parser

import (
	"encoding/json"

	"hms/gateway/pkg/docs/model"
)

func ParseComposition(inComposition []byte) (composition *model.Composition, err error) {
	composition = &model.Composition{}
	err = json.Unmarshal(inComposition, composition)
	return
}
