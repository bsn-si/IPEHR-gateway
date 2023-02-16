package treeindex

import (
	"fmt"
	"math/big"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/hm"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func DecryptDataValueNode(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	var err error

	switch node.Type {
	case base.DvURIItemType:
		err = decryptDvURI(node, key, nonce)
	case base.DvQuantityItemType:
		err = decryptDvQuantity(node, key, nonce)
	case base.DvTextItemType:
		err = decryptDvText(node, key, nonce)
	case base.DvCodedTextItemType:
		err = decryptDvCodedText(node, key, nonce)
	case base.DvDateTimeItemType:
		err = decryptDvDateTime(node, key)
	case base.DvIdentifierItemType:
		err = decryptDvIdentifier(node, key, nonce)
	case base.DvProportionItemType:
		err = decryptDvProportion(node, key)
	case base.DvMultimediaItemType:
		err = decryptDvMultimedia(node, key, nonce)
	default:
		fmt.Printf("decryptDataValueNode: unsupported node type: %s\n", node.Type)
	}

	return err
}

func decryptDvText(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	valueAttr, ok := node.Values["value"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvText.value")
	}

	newValue, err := hm.DecryptString(valueAttr.(*ValueNode).Data.([]byte), key, nonce)
	if err != nil {
		return fmt.Errorf("DecryptString error: %w", err)
	}

	node.Values["value"] = newValueNode(string(newValue))

	return nil
}

func decryptDvCodedText(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	valueAttr, ok := node.Values["value"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvCodedText.value")
	}

	newValue, err := hm.DecryptString(valueAttr.(*ValueNode).Data.([]byte), key, nonce)
	if err != nil {
		return fmt.Errorf("DecryptString error: %w", err)
	}

	node.Values["value"] = newValueNode(string(newValue))

	return nil
}

func decryptDvURI(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	valueAttr, ok := node.Values["value"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvURI.value")
	}

	newValue, err := hm.DecryptString(valueAttr.(*ValueNode).Data.([]byte), key, nonce)
	if err != nil {
		return fmt.Errorf("DecryptString error: %w", err)
	}

	node.Values["value"] = newValueNode(string(newValue))

	return nil
}

func decryptDvQuantity(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	magnitudeAttr, ok := node.Values["magnitude"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvQuantity.magnitude")
	}

	newValue, err := hm.DecryptFloat(magnitudeAttr.(*ValueNode).Data.(*big.Float), key)
	if err != nil {
		return fmt.Errorf("EncryptFloat error: %w %f", err, magnitudeAttr.(*ValueNode).Data.(*big.Float))
	}

	node.Values["magnitude"] = newValueNode(newValue)

	precisionAttr, ok := node.Values["precision"]
	if ok {
		newValue, err := hm.DecryptInt(precisionAttr.(*ValueNode).Data.(*big.Int), key)
		if err != nil {
			return fmt.Errorf("DecryptInt64 error: %w value: %d", err, precisionAttr.(*ValueNode).Data.(int64))
		}

		node.Values["precision"] = newValueNode(newValue)
	}

	unitsAttr, ok := node.Values["units"]
	if ok {
		newValue, err := hm.DecryptString(unitsAttr.(*ValueNode).Data.([]byte), key, nonce)
		if err != nil {
			return fmt.Errorf("DecryptString error: %w", err)
		}

		node.Values["units"] = newValueNode(string(newValue))
	}

	unitsSystemAttr, ok := node.Values["units_system"]
	if ok {
		newValue, err := hm.DecryptString(unitsSystemAttr.(*ValueNode).Data.([]byte), key, nonce)
		if err != nil {
			return fmt.Errorf("DecryptString error: %w", err)
		}

		node.Values["units_system"] = newValueNode(string(newValue))
	}

	unitsDisplayNameAttr, ok := node.Values["units_display_name"]
	if ok {
		newValue, err := hm.DecryptString(unitsDisplayNameAttr.(*ValueNode).Data.([]byte), key, nonce)
		if err != nil {
			return fmt.Errorf("DecryptString error: %w", err)
		}

		node.Values["units_display_name"] = newValueNode(string(newValue))
	}

	// todo DvAmount, DvQuantity.NormalRange, DvQuantity.OtherReferenceRanges

	return nil
}

func decryptDvDateTime(node *DataValueNode, key *hm.Key) error {
	valueAttr, ok := node.Values["value"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvDateTime.value")
	}

	value := valueAttr.(*ValueNode).Data.(*big.Int)

	newValue, err := hm.DecryptInt(value, key)
	if err != nil {
		return fmt.Errorf("DecryptInt64 error: %w value: %d", err, value)
	}

	node.Values["value"] = newValueNode(time.UnixMicro(newValue).UTC().Format(common.OpenEhrTimeFormatMicro))

	return nil
}

func decryptDvIdentifier(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	attributes := []string{"issuser", "assigner", "id", "type"}

	for _, attr := range attributes {
		valueAttr, ok := node.Values[attr]
		if !ok {
			return errors.ErrFieldIsEmpty("DvIdentifier." + attr)
		}

		newValue, err := hm.DecryptString(valueAttr.(*ValueNode).Data.([]byte), key, nonce)
		if err != nil {
			return fmt.Errorf("DecryptString error: %w", err)
		}

		node.Values[attr] = newValueNode(string(newValue))
	}

	return nil
}

func decryptDvProportion(node *DataValueNode, key *hm.Key) error {
	valueAttr, ok := node.Values["numerator"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvProportion.numerator")
	}

	switch value := valueAttr.(type) {
	case *ValueNode:
		switch f := value.Data.(type) {
		case *big.Float:
			newValue, err := hm.DecryptFloat(f, key)
			if err != nil {
				return fmt.Errorf("DecryptFloat64 error: %w value: %f", err, f)
			}

			node.Values["numerator"] = newValueNode(newValue)
		default:
			return errors.Errorf("decryptDvProportion: unsupported numerator.data type: %T", f)
		}
	default:
		return errors.Errorf("decryptDvProportion: unsupported numerator type: %T", value)
	}

	return nil
}

func decryptDvMultimedia(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	uriAttr, ok := node.Values["uri"]
	if ok {
		switch uri := uriAttr.(type) {
		case *DataValueNode:
			valueAttr, ok := uri.Values["value"]
			if !ok {
				return errors.ErrFieldIsEmpty("DvURI.value")
			}

			value := valueAttr.(*ValueNode).Data.([]byte)

			newValue, err := hm.DecryptString(value, key, nonce)
			if err != nil {
				return fmt.Errorf("DecryptString error: %w", err)
			}

			uri.Values["value"] = newValueNode(string(newValue))
		default:
			return errors.Errorf("decryptDvMultimedia: unsupported uri type: %T", uri)
		}
	}

	return nil
}
