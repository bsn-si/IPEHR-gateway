package treeindex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hms/gateway/pkg/docs/model/base"
)

type nodeType byte

const (
	objectNodeType nodeType = iota
	sliceNodeType
	dataValueNodeType
	valueNodeType
)

type noder interface {
	getNodeType() nodeType

	getID() string

	addAttribute(key string, val noder)
}

type baseNode struct {
	ID   string        `json:"id,omitempty"`
	Type base.ItemType `json:"type,omitempty"`
	Name string        `json:"name,omitempty"`
}

type objectNode struct {
	baseNode

	attributesOrder []string
	attributes      map[string]noder `json:"-"`
}

func (node objectNode) getNodeType() nodeType {
	return objectNodeType
}

func (node objectNode) getID() string {
	return node.ID
}

func (node *objectNode) addAttribute(key string, val noder) {
	node.attributesOrder = append(node.attributesOrder, key)
	node.attributes[key] = val
}

func (node objectNode) MarshalJSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	fmt.Fprintf(buffer, "{")
	fmt.Fprintf(buffer, `"id":"%s",`, node.ID)
	fmt.Fprintf(buffer, `"name":"%s",`, node.Name)
	fmt.Fprintf(buffer, `"type":"%s"`, node.Type)

	for _, k := range node.attributesOrder {
		data, err := json.Marshal(node.attributes[k])
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(buffer, `,"%s":%s`, k, string(data))
	}

	fmt.Fprintf(buffer, "}")
	return buffer.Bytes(), nil
}

type sliceNode struct {
	data map[string]noder
}

func (node sliceNode) getNodeType() nodeType {
	return sliceNodeType
}

func (node sliceNode) getID() string {
	return ""
}

func (node sliceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(node.data)
}

func (node *sliceNode) addAttribute(key string, val noder) {
	node.data[val.getID()] = val
}

type dataValueNode struct {
	baseNode
	Values map[string]noder `json:"values,omitempty"`
}

func (node dataValueNode) getNodeType() nodeType {
	return dataValueNodeType
}

func (node dataValueNode) getID() string {
	return node.ID
}

func (node *dataValueNode) addAttribute(key string, val noder) {
	node.Values[key] = val
}

type valueNode struct {
	baseNode
	data any
}

func (node valueNode) getNodeType() nodeType {
	return valueNodeType
}

func (node valueNode) getID() string {
	return ""
}

func (node *valueNode) addAttribute(key string, val noder) {
	noderInstance, ok := node.data.(noder)
	if !ok {
		return
	}

	noderInstance.addAttribute(key, val)
}

func (node valueNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(node.data)
}

func newNode(obj any) noder {
	switch obj := obj.(type) {
	case base.Root:
		return newObjectNode(obj)
	case base.DataValue:
		return newDataValueNode(obj)
	case base.CodePhrase:
		return nodeForCodePhrase(obj)
	default:
		return newValueNode(obj)
	}
}

func newObjectNode(obj base.Root) noder {
	l := obj.GetLocatable()

	return &objectNode{
		baseNode: baseNode{
			ID:   l.ArchetypeNodeID,
			Type: l.Type,
			Name: l.Name.Value,
		},
		attributes: map[string]noder{},
	}
}

func newSliceNode() noder {
	return &sliceNode{
		data: make(map[string]noder),
	}
}

func newDataValueNode(dv base.DataValue) noder {
	return &dataValueNode{
		baseNode: baseNode{
			ID:   "",
			Type: dv.GetType(),
		},
		Values: make(map[string]noder),
	}
}

func nodeForCodePhrase(cp base.CodePhrase) noder {
	return &valueNode{
		baseNode: baseNode{
			Type: cp.Type,
		},
		data: map[string]interface{}{
			"terminology_id": &valueNode{
				baseNode: baseNode{
					Type: cp.TerminologyID.Type,
				},
				data: map[string]interface{}{
					"value": cp.TerminologyID.Value,
				},
			},
			"code_string":    cp.CodeString,
			"preferred_term": cp.PreferredTerm,
		},
	}
}

func newValueNode(val any) noder {
	return &valueNode{
		data: val,
	}
}
