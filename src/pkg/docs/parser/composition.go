package parser

import (
	"encoding/json"

	"hms/gateway/pkg/docs/model"
)

func ParseComposition(inComposition []byte) (composition *model.Composition, err error) {
	composition = &model.Composition{}
	return json.Unmarshal(inComposition, composition)
}
