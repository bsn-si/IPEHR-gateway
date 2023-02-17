package treeindex

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func ExecNode(node Noder, dataValueExec func(node *DataValueNode) error) (Noder, error) {
	switch node := node.(type) {
	case *CompositionNode:
		if err := execCompositionData(node, dataValueExec); err != nil {
			return nil, fmt.Errorf("encryptCompositionData error: %w", err)
		}
	default:
		fmt.Printf("EncryptNode: Unsupported type %T", node)
	}

	return node, nil
}

func execCompositionData(node *CompositionNode, dataValueExec func(node *DataValueNode) error) error {
	for _, container := range node.Data {
		for _, collection := range container {
			for _, node := range collection {
				switch node.GetNodeType() {
				case NodeTypeObject:
					err := execObjectNode(node.(*ObjectNode), dataValueExec)
					if err != nil {
						return fmt.Errorf("encryptObjectNode error: %w", err)
					}
				default:
					return errors.Errorf("Unsupported collection node type: %d", node.GetNodeType())
				}
			}
		}
	}

	return nil
}

func execObjectNode(nodeObj *ObjectNode, dataValueExec func(node *DataValueNode) error) error {
	switch nodeObj.Type {
	case base.InstructionItemType:
		if err := encryptInstruction(nodeObj); err != nil {
			return fmt.Errorf("encryptInstruction error: %w", err)
		}
	case base.ActionItemType:
		if err := encryptAction(nodeObj); err != nil {
			return fmt.Errorf("encrypt ACTION error: %w", err)
		}
	case base.EvaluationItemType:
		if err := execByAttribute(nodeObj, "data", true, dataValueExec); err != nil {
			return fmt.Errorf("encrypt EVALUATION error: %w", err)
		}
	case base.ObservationItemType:
		if err := execByAttribute(nodeObj, "data", true, dataValueExec); err != nil {
			return fmt.Errorf("encrypt OBSERVATION error: %w", err)
		}

		if err := execByAttribute(nodeObj, "protocol", false, dataValueExec); err != nil {
			return fmt.Errorf("encrypt OBSERVATION error: %w", err)
		}
	case base.HistoryItemType:
		if err := execByAttribute(nodeObj, "events", true, dataValueExec); err != nil {
			return fmt.Errorf("encrypt HISTORY error: %w", err)
		}
	case base.ElementItemType:
		if err := execByAttribute(nodeObj, "value", true, dataValueExec); err != nil {
			return fmt.Errorf("encrypt ELEMENT error: %w", err)
		}
	case base.PointEventItemType:
		if err := execByAttribute(nodeObj, "data", true, dataValueExec); err != nil {
			return fmt.Errorf("encrypt POINT_EVENT error: %w", err)
		}
	case base.ItemTreeItemType:
		if err := execByAttribute(nodeObj, "items", true, dataValueExec); err != nil {
			return fmt.Errorf("encrypt ITEM_TREE error: %w", err)
		}
	case base.ClusterItemType:
		if err := execByAttribute(nodeObj, "items", true, dataValueExec); err != nil {
			return fmt.Errorf("encrypt CLUSTER error: %w", err)
		}
	default:
		fmt.Printf("execObjectNode: unsupported ObjectNode.Type: %s\n", nodeObj.Type)
	}

	return nil
}

func execByAttribute(nodeObj *ObjectNode, attribute string, reqired bool, dataValueExec func(node *DataValueNode) error) error {
	attr, ok := nodeObj.Attributes[attribute]
	if !ok && reqired {
		return errors.ErrFieldIsEmpty("Attributes." + attribute)
	} else if !ok && !reqired {
		return nil
	}

	switch attr := attr.(type) {
	case *ObjectNode:
		if err := execObjectNode(attr, dataValueExec); err != nil {
			return fmt.Errorf("execObjectNode error: %w", err)
		}
	case *SliceNode:
		if err := encryptWalkBySlice(attr, dataValueExec); err != nil {
			return fmt.Errorf("encryptWalkBySlice error: %w", err)
		}
	case *DataValueNode:
		if err := dataValueExec(attr); err != nil {
			return fmt.Errorf("dataValueExec error: %w", err)
		}
	default:
		fmt.Printf("execByAttribute: unsupported attribute type: %T attr: %s\n", attr, attribute)
	}

	return nil
}

func encryptInstruction(nodeObj *ObjectNode) error {
	//TODO
	return nil
}

func encryptAction(nodeObj *ObjectNode) error {
	//TODO
	return nil
}

func encryptWalkBySlice(slice *SliceNode, dataValueExec func(node *DataValueNode) error) error {
	for _, node := range slice.Data {
		switch node := node.(type) {
		case *ObjectNode:
			if err := execObjectNode(node, dataValueExec); err != nil {
				return fmt.Errorf("encryptObjectNode error: %w", err)
			}
		default:
			fmt.Printf("encryptWalkBySlice: unsupported attribute type: %T\n", node)
		}
	}

	return nil
}
