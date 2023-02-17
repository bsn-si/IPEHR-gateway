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
	case NodeTypeObject:
		nw.data = &ObjectNode{}
	case NodeTypeSlice:
		nw.data = newSliceNode()
	case NodeTypeDataValue:
		nw.data = &DataValueNode{}
	case NodeTypeValue:
		nw.data = &ValueNode{}
	case NodeTypeEHR:
		nw.data = &EHRNode{}
	case NodeTypeCompostion:
		nw.data = &CompositionNode{}
	case NodeTypeEventContext:
		nw.data = &EventContextNode{}
	default:
		return fmt.Errorf("unexpected node type: %v", tmp.NodeType) //nolint
	}

	if err := msgpack.Unmarshal(data, &nw.data); err != nil {
		return errors.Wrap(err, "cannot decode node")
	}

	return nil
}
