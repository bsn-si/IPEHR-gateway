package treeindex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
)

func processEHR(ehr model.EHR) (*EHRNode, error) {
	node := newEHRNode(ehr)

	for _, cmp := range ehr.Compositions {
		if err := node.addComposition(cmp); err != nil {
			return nil, err
		}
	}

	node.addAttribute("system_id", newNode(ehr.SystemID))
	node.addAttribute("ehr_id", newNode(ehr.EhrID))

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

func (ehr *EHRNode) addComposition(cmp model.Composition) error {
	cmpNode, err := processComposition(cmp)
	if err != nil {
		return errors.Wrap(err, "cannot add Composition node into EHRNode")
	}

	ehr.compositions[cmpNode.GetID()] = append(ehr.compositions[cmpNode.GetID()], cmpNode)
	return nil
}

func (ehr EHRNode) GetCompositions() Container {
	return ehr.compositions
}

func (ehr EHRNode) GetNodeType() NodeType {
	return EHRNodeType
}

func (ehr EHRNode) GetID() string {
	return ehr.ID
}

func (ehr EHRNode) TryGetChild(key string) Noder {
	n := ehr.baseNode.TryGetChild(key)
	if n != nil {
		return n
	}

	return ehr.attributes[key]
}

func (ehr EHRNode) ForEach(f func(name string, node Noder) bool) {
	for key, node := range ehr.attributes {
		if !f(key, node) {
			break
		}
	}
}

func (ehr EHRNode) addAttribute(key string, val Noder) {
	ehr.attributes[key] = val
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
