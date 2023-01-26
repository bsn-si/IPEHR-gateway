package treeindex

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
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
	BaseNode

	Attributes   Attributes `json:"-"`
	Compositions Container
}

func newEHRNode(ehr model.EHR) *EHRNode {
	node := EHRNode{
		BaseNode: BaseNode{
			ID:       ehr.EhrID.Value,
			Type:     base.EHRItemType,
			NodeType: EHRNodeType,
		},
		Attributes:   Attributes{},
		Compositions: Container{},
	}

	return &node
}

func (ehr *EHRNode) addComposition(cmp model.Composition) error {
	cmpNode, err := processComposition(cmp)
	if err != nil {
		return errors.Wrap(err, "cannot add Composition node into EHRNode")
	}

	ehr.Compositions[cmpNode.GetID()] = append(ehr.Compositions[cmpNode.GetID()], cmpNode)
	return nil
}

func (ehr EHRNode) GetCompositions() Container {
	return ehr.Compositions
}

func (ehr EHRNode) GetID() string {
	return ehr.ID
}

func (ehr EHRNode) TryGetChild(key string) Noder {
	n := ehr.BaseNode.TryGetChild(key)
	if n != nil {
		return n
	}

	return ehr.Attributes[key]
}

func (ehr EHRNode) addAttribute(key string, val Noder) {
	ehr.Attributes[key] = val
}

func (ehr EHRNode) MarshalJSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	fmt.Fprintf(buffer, "{")
	fmt.Fprintf(buffer, `"id":"%s",`, ehr.ID)
	fmt.Fprintf(buffer, `"name":"%s",`, ehr.Name)
	fmt.Fprintf(buffer, `"type":"%s"`, ehr.Type)

	if l := ehr.Compositions.Len(); l > 0 {
		cmps := make([]Noder, 0, l)
		for _, nodes := range ehr.Compositions {
			cmps = append(cmps, nodes...)
		}

		cmpsData, err := json.Marshal(cmps)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(buffer, `,"compositions":%s`, string(cmpsData))
	}

	for k, v := range ehr.Attributes {
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(buffer, `,"%s":%s`, k, string(data))
	}

	fmt.Fprintf(buffer, "}")
	return buffer.Bytes(), nil
}
