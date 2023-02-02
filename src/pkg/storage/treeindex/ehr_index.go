package treeindex

import (
	"encoding/json"
	"fmt"

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
	node, err := ProcessEHR(ehr)
	if err != nil {
		return errors.Wrap(err, "cannot process EHR")
	}

	err = idx.AddEHRNode(node)
	if err != nil {
		return errors.Wrap(err, "cannot add EHR node")
	}

	return nil
}

func (idx *EHRIndex) AddEHRNode(node *EHRNode) error {
	if node == nil {
		return errors.New("node is empty")
	}

	nodeID := node.GetID()
	if nodeID == "" {
		return errors.New("nodeID is empty")
	}

	if idx.Ehrs == nil {
		idx.Ehrs = map[string]*EHRNode{}
	}

	if _, ok := idx.Ehrs[nodeID]; ok {
		return errors.New("EHR nodeID already exists")
	}

	idx.Ehrs[nodeID] = node

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

	if err := ehrNode.addComposition(cmp); err != nil {
		return fmt.Errorf("ehrNode.addComposition error: %w", err)
	}

	return nil
}

func (idx EHRIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(idx.Ehrs)
}
