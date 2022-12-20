package treeindex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
)

func processEHR(ehr model.EHR) (Noder, error) {
	node := newEHRNode(ehr)

	for _, cmp := range ehr.Compositions {
		cmpNode, err := walk(cmp)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get node for Composition object")
		}

		node.compositions[cmpNode.GetID()] = append(node.compositions[cmpNode.GetID()], cmpNode)
	}

	// TODO: add logic for other fields
	return node, nil
}

type EHRNode struct {
	baseNode

	attributes   map[string]Noder `json:"-"`
	compositions Container
}

func newEHRNode(ehr model.EHR) *EHRNode {
	node := EHRNode{
		baseNode: baseNode{
			ID:   ehr.EhrID.Value,
			Type: base.EHRItemType,
		},
		attributes:   map[string]Noder{},
		compositions: Container{},
	}

	return &node
}

func (ehr EHRNode) GetNodeType() NodeType {
	return EHRNodeType
}

func (ehr EHRNode) GetID() string {
	return ehr.ID
}

func (ehr EHRNode) TryGetChild(key string) Noder {
	return nil
}

func (ehr EHRNode) ForEach(func(name string, node Noder) bool) {
}

func (ehr EHRNode) addAttribute(key string, val Noder) {

}

func (ehr EHRNode) MarshalJSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	fmt.Fprintf(buffer, "{")
	fmt.Fprintf(buffer, `"id":"%s",`, ehr.ID)
	fmt.Fprintf(buffer, `"name":"%s",`, ehr.Name)
	fmt.Fprintf(buffer, `"type":"%s"`, ehr.Type)

	if l := ehr.compositions.Len(); l > 0 {
		cmps := make([]Noder, 0, l)
		for _, nodes := range ehr.compositions {
			cmps = append(cmps, nodes...)
		}

		cmpsData, err := json.Marshal(cmps)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(buffer, `,"compositions":%s`, string(cmpsData))
	}

	for k, v := range ehr.attributes {
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(buffer, `,"%s":%s`, k, string(data))
	}

	fmt.Fprintf(buffer, "}")
	return buffer.Bytes(), nil
}
