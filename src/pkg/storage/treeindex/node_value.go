package treeindex

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

type DataType = uint8

const (
	DataTypeBasic DataType = iota
	DataTypeBigFloat
	DataTypeBigInt
)

type ValueNode struct {
	BaseNode
	DataType DataType `msgpack:"DataType,omitempty"`
	Data     any
}

func newValueNode(val any) Noder {
	var dt DataType

	switch val.(type) {
	case *big.Int:
		dt = DataTypeBigInt
	case *big.Float:
		dt = DataTypeBigFloat
	default:
		dt = DataTypeBasic
	}

	return &ValueNode{
		BaseNode: BaseNode{
			NodeType: NodeTypeValue,
		},
		DataType: dt,
		Data:     val,
	}
}

func (node ValueNode) GetData() any {
	return node.Data
}

func (node ValueNode) GetID() string {
	return ""
}

func (node ValueNode) TryGetChild(key string) Noder {
	n := node.BaseNode.TryGetChild(key)
	if n != nil {
		return n
	}

	return nil
}

func (node *ValueNode) addAttribute(key string, val Noder) {
	noderInstance, ok := node.Data.(Noder)
	if !ok {
		return
	}

	noderInstance.addAttribute(key, val)
}

func (node ValueNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(node.Data)
}

func (node *ValueNode) UnmarshalMsgpack(data []byte) error {
	tmp := struct {
		BaseNode
		DataType DataType
		Data     any
	}{}

	if err := msgpack.Unmarshal(data, &tmp); err != nil {
		return err
	}

	node.BaseNode = tmp.BaseNode
	node.DataType = tmp.DataType

	switch tmp.DataType {
	case DataTypeBasic:
		switch v := tmp.Data.(type) {
		case int8:
			node.Data = int(v)
		case int16:
			node.Data = int(v)
		case uint16:
			node.Data = int(v)
		case int32:
			node.Data = int(v)
		case uint32:
			node.Data = int(v)
		default:
			node.Data = tmp.Data
		}
	case DataTypeBigInt:
		switch v := tmp.Data.(type) {
		case []uint8:
			node.Data = new(big.Int).SetBytes(tmp.Data.([]uint8))
		default:
			return errors.Errorf("Unsupported ValueNode type for DataTypeBigInt. Expected []uint8, received: %T", v)
		}
	case DataTypeBigFloat:
		switch v := tmp.Data.(type) {
		case []uint8:
			node.Data = big.NewFloat(0)

			err := node.Data.(*big.Float).UnmarshalText(v)
			if err != nil {
				return fmt.Errorf("big.Float UnmarshalText error: %w", err)
			}
		default:
			return errors.Errorf("Unsupported ValueNode type for DataTypeBigFloat. Expected []uint8, received: %T", v)
		}
	default:
		return errors.Errorf("Unsupported ValueNode DataType: %d", tmp.DataType)
	}

	return nil
}
