package treeindex

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
)

var DefaultEHRIndex = NewEHRIndex()

type EHRIndex struct {
	ehrs map[string]*EHRNode
}

func NewEHRIndex() *EHRIndex {
	idx := EHRIndex{
		ehrs: map[string]*EHRNode{},
	}

	return &idx
}

func (idx *EHRIndex) AddEHR(ehr model.EHR) error {
	node, err := processEHR(ehr)

	if err != nil {
		return errors.Wrap(err, "cannot add EHR object")
	}

	idx.ehrs[node.GetID()] = node

	return nil
}

func (idx *EHRIndex) AddComposition(ehrID string, cmp model.Composition) error {
	ehrNode, ok := idx.ehrs[ehrID]
	if !ok {
		return errors.New("EHR not found")
	}

	return ehrNode.addComposition(cmp)
}

func (idx EHRIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(idx.ehrs)
}
