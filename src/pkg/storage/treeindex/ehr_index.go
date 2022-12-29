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

func AddEHR(ehr model.EHR) error {
	return DefaultEHRIndex.AddEHR(ehr)
}

func AddComposition(ehrID string, cmp model.Composition) error {
	return DefaultEHRIndex.AddComposition(ehrID, cmp)
}

func (idx *EHRIndex) AddEHR(ehr model.EHR) error {
	node, err := processEHR(ehr)

	if err != nil {
		return errors.Wrap(err, "cannot add EHR object")
	}

	idx.ehrs[node.GetID()] = node

	return nil
}

func (idx EHRIndex) GetEHRs(id string) ([]*EHRNode, error) {
	if id == "" {
		result := make([]*EHRNode, 0, len(idx.ehrs))
		for _, ehrNode := range idx.ehrs {
			result = append(result, ehrNode)
		}

		return result, nil
	}

	ehrNode, ok := idx.ehrs[id]
	if !ok {
		return nil, errors.New("cannot get EHR by id")
	}

	return []*EHRNode{ehrNode}, nil
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
