package treeindex

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/hm"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func EncryptNode(node Noder, key *hm.Key, nonce *hm.Nonce) (Noder, error) {
	switch node := node.(type) {
	case *CompositionNode:
		if err := encryptCompositionData(node, key, nonce); err != nil {
			return nil, fmt.Errorf("encryptCompositionData error: %w", err)
		}
	default:
		fmt.Printf("EncryptNode: Unsupported type %T", node)
	}

	return node, nil
}

func encryptCompositionData(node *CompositionNode, key *hm.Key, nonce *hm.Nonce) error {
	for _, container := range node.Data {
		for _, collection := range container {
			for _, node := range collection {
				switch node.GetNodeType() {
				case ObjectNodeType:
					err := encryptObjectNode(node.(*ObjectNode), key, nonce)
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

func encryptObjectNode(nodeObj *ObjectNode, key *hm.Key, nonce *hm.Nonce) error {
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
		if err := encryptByAttribute(nodeObj, "data", true, key, nonce); err != nil {
			return fmt.Errorf("encrypt EVALUATION error: %w", err)
		}
	case base.ObservationItemType:
		if err := encryptByAttribute(nodeObj, "data", true, key, nonce); err != nil {
			return fmt.Errorf("encrypt OBSERVATION error: %w", err)
		}

		if err := encryptByAttribute(nodeObj, "protocol", false, key, nonce); err != nil {
			return fmt.Errorf("encrypt OBSERVATION error: %w", err)
		}
	case base.HistoryItemType:
		if err := encryptByAttribute(nodeObj, "events", true, key, nonce); err != nil {
			return fmt.Errorf("encrypt HISTORY error: %w", err)
		}
	case base.ElementItemType:
		if err := encryptByAttribute(nodeObj, "value", true, key, nonce); err != nil {
			return fmt.Errorf("encrypt ELEMENT error: %w", err)
		}
	case base.PointEventItemType:
		if err := encryptByAttribute(nodeObj, "data", true, key, nonce); err != nil {
			return fmt.Errorf("encrypt POINT_EVENT error: %w", err)
		}
	case base.ItemTreeItemType:
		if err := encryptByAttribute(nodeObj, "items", true, key, nonce); err != nil {
			return fmt.Errorf("encrypt ITEM_TREE error: %w", err)
		}
	case base.ClusterItemType:
		if err := encryptByAttribute(nodeObj, "items", true, key, nonce); err != nil {
			return fmt.Errorf("encrypt CLUSTER error: %w", err)
		}
	default:
		fmt.Printf("encryptObjectNode: unsupported ObjectNode.Type: %s\n", nodeObj.Type)
	}

	return nil
}

func encryptByAttribute(nodeObj *ObjectNode, attribute string, reqired bool, key *hm.Key, nonce *hm.Nonce) error {
	attr, ok := nodeObj.Attributes[attribute]
	if !ok && reqired {
		return errors.ErrFieldIsEmpty("Attributes." + attribute)
	}

	switch attr := attr.(type) {
	case *ObjectNode:
		if err := encryptObjectNode(attr, key, nonce); err != nil {
			return fmt.Errorf("encryptObjectNode error: %w", err)
		}
	case *SliceNode:
		if err := encryptWalkBySlice(attr, key, nonce); err != nil {
			return fmt.Errorf("encryptWalkBySlice error: %w", err)
		}
	case *DataValueNode:
		if err := encryptDataValueNode(attr, key, nonce); err != nil {
			return fmt.Errorf("encryptDataValueNode error: %w", err)
		}
	default:
		fmt.Printf("encryptByAttribute: unsupported attribute type: %T\n", attr)
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

func encryptWalkBySlice(slice *SliceNode, key *hm.Key, nonce *hm.Nonce) error {
	for _, node := range slice.Data {
		switch node := node.(type) {
		case *ObjectNode:
			if err := encryptObjectNode(node, key, nonce); err != nil {
				return fmt.Errorf("encryptObjectNode error: %w", err)
			}
		default:
			fmt.Printf("TODO: encryptWalkBySlice: %T\n", node)
		}
	}

	return nil
}

func encryptDataValueNode(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	var err error

	switch node.Type {
	case base.DvURIItemType:
		err = encryptDvURI(node, key, nonce)
	case base.DvQuantityItemType:
		err = encryptDvQuantity(node, key, nonce)
	case base.DvTextItemType:
		err = encryptDvText(node, key, nonce)
	case base.DvCodedTextItemType:
		err = encryptDvCodedText(node, key, nonce)
	case base.DvDateTimeItemType:
		err = encryptDvDateTime(node, key)
	case base.DvIdentifierItemType:
		err = encryptDvIdentifier(node, key, nonce)
	case base.DvProportionItemType:
		err = encryptDvProportion(node, key)
	case base.DvMultimediaItemType:
		err = encryptDvMultimedia(node, key, nonce)
	default:
		fmt.Printf("TODO: encryptDataValueNode: %s\n", node.Type)
	}

	return err
}
