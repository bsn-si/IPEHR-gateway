package treeindex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
)

type NodeType byte

const (
	ObjectNodeType NodeType = iota
	SliceNodeType
	DataValueNodeType
	ValueNodeType
	EHRNodeType
	CompostionNodeType
	EventContextNodeType
)

type Noder interface {
	GetNodeType() NodeType

	GetID() string

	addAttribute(key string, val Noder)
	TryGetChild(key string) Noder
	ForEach(func(name string, node Noder) bool)
}

type baseNode struct {
	ID   string        `json:"id,omitempty"`
	Type base.ItemType `json:"type,omitempty"`
	Name string        `json:"name,omitempty"`
}

type ObjectNode struct {
	baseNode

	attributesOrder []string
	attributes      map[string]Noder `json:"-"`
}

func (node ObjectNode) GetNodeType() NodeType {
	return ObjectNodeType
}

func (node ObjectNode) GetID() string {
	return node.ID
}

func (node ObjectNode) TryGetChild(key string) Noder {
	n, ok := node.attributes[key]
	if !ok {
		return nil
	}

	return n
}

func (node ObjectNode) ForEach(foo func(name string, node Noder) bool) {
	for key, node := range node.attributes {
		if !foo(key, node) {
			break
		}
	}
}

func (node *ObjectNode) addAttribute(key string, val Noder) {
	node.attributesOrder = append(node.attributesOrder, key)
	node.attributes[key] = val
}

func (node ObjectNode) MarshalJSON() ([]byte, error) {
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

type SliceNode struct {
	data map[string]Noder
}

func (node SliceNode) GetNodeType() NodeType {
	return SliceNodeType
}

func (node SliceNode) GetID() string {
	return ""
}

func (node SliceNode) TryGetChild(key string) Noder {
	n, ok := node.data[key]
	if !ok {
		return nil
	}

	return n
}

func (node SliceNode) ForEach(foo func(name string, node Noder) bool) {
	for key, node := range node.data {
		if !foo(key, node) {
			break
		}
	}
}

func (node SliceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(node.data)
}

func (node *SliceNode) addAttribute(key string, val Noder) {
	node.data[val.GetID()] = val
}

type DataValueNode struct {
	baseNode
	Values map[string]Noder `json:"values,omitempty"`
}

func (node DataValueNode) GetNodeType() NodeType {
	return DataValueNodeType
}

func (node DataValueNode) GetID() string {
	return node.ID
}

func (node DataValueNode) TryGetChild(key string) Noder {
	n, ok := node.Values[key]
	if !ok {
		return nil
	}

	return n
}

func (node DataValueNode) ForEach(foo func(name string, node Noder) bool) {
	for key, node := range node.Values {
		if !foo(key, node) {
			break
		}
	}
}

func (node *DataValueNode) addAttribute(key string, val Noder) {
	node.Values[key] = val
}

type ValueNode struct {
	baseNode
	data any
}

func (node ValueNode) GetData() any {
	return node.data
}

func (node ValueNode) GetNodeType() NodeType {
	return ValueNodeType
}

func (node ValueNode) GetID() string {
	return ""
}

func (node ValueNode) TryGetChild(key string) Noder {
	return nil
}

func (node ValueNode) ForEach(foo func(name string, node Noder) bool) {
}

func (node *ValueNode) addAttribute(key string, val Noder) {
	noderInstance, ok := node.data.(Noder)
	if !ok {
		return
	}

	noderInstance.addAttribute(key, val)
}

func (node ValueNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(node.data)
}

func newNode(obj any) Noder {
	switch obj := obj.(type) {
	case model.EHR:
		return newEHRNode(obj)
	case model.Composition:
		return newCompositionNode(obj)
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

func newObjectNode(obj base.Root) Noder {
	l := obj.GetLocatable()

	return &ObjectNode{
		baseNode: baseNode{
			ID:   l.ArchetypeNodeID,
			Type: l.Type,
			Name: l.Name.Value,
		},
		attributes: map[string]Noder{},
	}
}

func newSliceNode() Noder {
	return &SliceNode{
		data: make(map[string]Noder),
	}
}

func newDataValueNode(dv base.DataValue) Noder {
	return &DataValueNode{
		baseNode: baseNode{
			ID:   "",
			Type: dv.GetType(),
		},
		Values: make(map[string]Noder),
	}
}

func nodeForCodePhrase(cp base.CodePhrase) Noder {
	return &ValueNode{
		baseNode: baseNode{
			Type: cp.Type,
		},
		data: map[string]interface{}{
			"terminology_id": &ValueNode{
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

func newValueNode(val any) Noder {
	return &ValueNode{
		data: val,
	}
}
