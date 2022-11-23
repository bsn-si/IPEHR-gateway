package treeindex

import (
	"encoding/json"
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
	ID   string
	Type base.ItemType
	Name string
}

type objectNode struct {
	baseNode

	Attributes Attributes
	Value      map[string]interface{}
}

func (node objectNode) getNodeType() nodeType {
	return objectNodeType
}

func (node objectNode) getID() string {
	return node.ID
}

func (node *objectNode) addAttribute(key string, val noder) {
	node.Attributes.add(key, val)
}

type sliceNode struct {
	Data []noder
}

func (node sliceNode) getNodeType() nodeType {
	return sliceNodeType
}

type dataValueNode struct {
	baseNode
	Values map[string]noder
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

func NewNode(obj base.Root) noder {
	l := obj.GetLocatable()

	return &objectNode{
		baseNode: baseNode{
			ID:   l.ArchetypeNodeID,
			Type: l.Type,
			Name: l.Name.Value,
		},
		Attributes: map[string]map[string]noder{},
	}
}

func NewNodeForCodePhrase(mt base.CodePhrase) noder {
	return &objectNode{
		baseNode: baseNode{
			Type: mt.Type,
		},
		Value: map[string]interface{}{
			"terminology_id": &objectNode{
				baseNode: baseNode{
					Type: mt.TerminologyID.Type,
				},
				Value: map[string]interface{}{
					"value": mt.TerminologyID.Value,
				},
			},
			"code_string":    mt.CodeString,
			"preferred_term": mt.PreferredTerm,
		},
	}
}

func NewNodeForData(dv base.DataValue) noder {
	return &dataValueNode{
		baseNode: baseNode{
			ID:   "",
			Type: dv.GetType(),
		},
		Values: make(map[string]noder),
	}
}

func NewValueNode(val any) noder {
	return &valueNode{
		data: val,
	}
}

type Attributes map[string]map[string]noder

func (a Attributes) add(name string, node noder) {
	m, ok := a[name]
	if !ok {
		m = map[string]noder{}
		a[name] = m
	}

	m[node.getID()] = node
}
