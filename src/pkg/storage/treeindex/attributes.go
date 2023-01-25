package treeindex

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

type Attributes map[string]Noder

func (attributes *Attributes) DecodeMsgpack(dec *msgpack.Decoder) error {
	tempMap := map[string]NodeWrapper{}
	if err := dec.Decode(&tempMap); err != nil {
		return errors.Wrap(err, "cannot unmarshal attributes map")
	}

	attr := Attributes{}
	for k, v := range tempMap {
		attr[k] = v.data
	}

	*attributes = attr
	return nil
}

type NodeWrapper struct {
	data Noder `msgpack:"-"`
}

func (nw *NodeWrapper) UnmarshalMsgpack(data []byte) error {
	tmp := struct {
		NodeType NodeType `msgpack:"node_type"`
	}{}

	if err := msgpack.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "cannot decode tmp struct")
	}

	switch tmp.NodeType {
	case ObjectNodeType:
		nw.data = &ObjectNode{}
	case SliceNodeType:
		nw.data = newSliceNode()
	case DataValueNodeType:
		nw.data = &DataValueNode{}
	case ValueNodeType:
		nw.data = &ValueNode{}
	case EHRNodeType:
		nw.data = &EHRNode{}
	case CompostionNodeType:
		nw.data = &CompositionNode{}
	case EventContextNodeType:
		nw.data = &EventContextNode{}
	default:
		return fmt.Errorf("unexpected node type: %v", tmp.NodeType) //nolint
	}

	if err := msgpack.Unmarshal(data, &nw.data); err != nil {
		return errors.Wrap(err, "cannot decode node")
	}

	return nil
}
