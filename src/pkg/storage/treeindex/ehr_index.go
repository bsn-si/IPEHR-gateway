package treeindex

import (
	"encoding/json"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

var DefaultEHRIndex = NewEHRIndex()

type EHRIndex struct {
	Ehrs map[string]*EHRNode `msgpack:"ehr,omitempty"`
}

func NewEHRIndex() *EHRIndex {
	idx := EHRIndex{
		Ehrs: map[string]*EHRNode{},
	}

	return &idx
}

func AddEHR(ehr *model.EHR) error {
	return DefaultEHRIndex.AddEHR(ehr)
}

func AddComposition(ehrID string, cmp *model.Composition) error {
	return DefaultEHRIndex.AddComposition(ehrID, cmp)
}

func (idx *EHRIndex) AddEHR(ehr *model.EHR) error {
	node, err := processEHR(ehr)
	if err != nil {
		return errors.Wrap(err, "cannot add EHR object")
	}

	idx.Ehrs[node.GetID()] = node

	return nil
}

func (idx EHRIndex) GetEHRs(id string) ([]*EHRNode, error) {
	if id == "" {
		result := make([]*EHRNode, 0, len(idx.Ehrs))
		for _, ehrNode := range idx.Ehrs {
			result = append(result, ehrNode)
		}

		return result, nil
	}

	ehrNode, ok := idx.Ehrs[id]
	if !ok {
		return nil, errors.New("cannot get EHR by id")
	}

	return []*EHRNode{ehrNode}, nil
}

func (idx *EHRIndex) AddComposition(ehrID string, cmp *model.Composition) error {
	ehrNode, ok := idx.Ehrs[ehrID]
	if !ok {
		return errors.New("EHR not found")
	}

	return ehrNode.addComposition(cmp)
}

func (idx EHRIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(idx.Ehrs)
}
