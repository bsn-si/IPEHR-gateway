package treeindex

import (
	"fmt"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/hm"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

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
		fmt.Printf("encryptDataValueNode: unsupported node type: %s\n", node.Type)
	}

	return err
}

func encryptDvText(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	valueAttr, ok := node.Values["value"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvText.value")
	}

	newValue := hm.EncryptString(valueAttr.(*ValueNode).Data.(string), key, nonce)

	node.Values["value"] = newValueNode(newValue)

	return nil
}

func encryptDvCodedText(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	valueAttr, ok := node.Values["value"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvCodedText.value")
	}

	newValue := hm.EncryptString(valueAttr.(*ValueNode).Data.(string), key, nonce)

	node.Values["value"] = newValueNode(newValue)

	return nil
}

func encryptDvURI(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	valueAttr, ok := node.Values["value"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvURI.value")
	}

	newValue := hm.EncryptString(valueAttr.(*ValueNode).Data.(string), key, nonce)

	node.Values["value"] = newValueNode(newValue)

	return nil
}

func encryptDvQuantity(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	magnitudeAttr, ok := node.Values["magnitude"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvQuantity.magnitude")
	}

	newValue, err := hm.EncryptFloat64(magnitudeAttr.(*ValueNode).Data.(float64), key)
	if err != nil {
		return fmt.Errorf("EncryptFloat error: %w %f", err, magnitudeAttr.(*ValueNode).Data.(float64))
	}

	node.Values["magnitude"] = newValueNode(newValue)

	precisionAttr, ok := node.Values["precision"]
	if ok {
		newValue, err := hm.EncryptInt(precisionAttr.(*ValueNode).Data.(int64), key)
		if err != nil {
			return fmt.Errorf("EncryptInt error: %w value: %d", err, precisionAttr.(*ValueNode).Data.(int))
		}

		node.Values["precision"] = newValueNode(newValue)
	}

	unitsAttr, ok := node.Values["units"]
	if ok {
		newValue := hm.EncryptString(unitsAttr.(*ValueNode).Data.(string), key, nonce)
		node.Values["units"] = newValueNode(newValue)
	}

	unitsSystemAttr, ok := node.Values["units_system"]
	if ok {
		newValue := hm.EncryptString(unitsSystemAttr.(*ValueNode).Data.(string), key, nonce)
		node.Values["units_system"] = newValueNode(newValue)
	}

	unitsDisplayNameAttr, ok := node.Values["units_display_name"]
	if ok {
		newValue := hm.EncryptString(unitsDisplayNameAttr.(*ValueNode).Data.(string), key, nonce)
		node.Values["units_display_name"] = newValueNode(newValue)
	}

	// todo DvAmount, DvQuantity.NormalRange, DvQuantity.OtherReferenceRanges

	return nil
}

func encryptDvDateTime(node *DataValueNode, key *hm.Key) error {
	valueAttr, ok := node.Values["value"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvDateTime.value")
	}

	value, err := time.Parse(common.OpenEhrTimeFormatMicro, valueAttr.(*ValueNode).Data.(string))
	if err != nil {
		return fmt.Errorf("time.Parse error: %w value: %s", err, valueAttr.(*ValueNode).Data.(string))
	}

	newValue, err := hm.EncryptInt(value.UnixMicro(), key)
	if err != nil {
		return fmt.Errorf("EncryptInt64 error: %w value: %d", err, value.Unix())
	}

	node.Values["value"] = newValueNode(newValue)

	return nil
}

func encryptDvIdentifier(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	attributes := []string{"issuser", "assigner", "id", "type"}

	for _, attr := range attributes {
		valueAttr, ok := node.Values[attr]
		if !ok {
			return errors.ErrFieldIsEmpty("DvIdentifier." + attr)
		}

		newValue := hm.EncryptString(valueAttr.(*ValueNode).Data.(string), key, nonce)
		node.Values[attr] = newValueNode(newValue)
	}

	return nil
}

func encryptDvProportion(node *DataValueNode, key *hm.Key) error {
	valueAttr, ok := node.Values["numerator"]
	if !ok {
		return errors.ErrFieldIsEmpty("DvProportion.numerator")
	}

	switch value := valueAttr.(type) {
	case *ValueNode:
		switch f := value.Data.(type) {
		case float64:
			newValue, err := hm.EncryptFloat64(f, key)
			if err != nil {
				return fmt.Errorf("EncryptFloat64 error: %w value: %f", err, f)
			}

			node.Values["numerator"] = newValueNode(newValue)
		default:
			return errors.Errorf("encryptDvProportion: unsupported numerator.data type: %T", f)
		}
	default:
		return errors.Errorf("encryptDvProportion: unsupported numerator type: %T", value)
	}

	return nil
}

func encryptDvMultimedia(node *DataValueNode, key *hm.Key, nonce *hm.Nonce) error {
	uriAttr, ok := node.Values["uri"]
	if ok {
		switch uri := uriAttr.(type) {
		case *DataValueNode:
			valueAttr, ok := uri.Values["value"]
			if !ok {
				return errors.ErrFieldIsEmpty("DvURI.value")
			}

			value := valueAttr.(*ValueNode).Data.(string)
			newValue := hm.EncryptString(value, key, nonce)
			uri.Values["value"] = newValueNode(newValue)
		default:
			return errors.Errorf("encryptDvMultimedia: unsupported uri type: %T", uri)
		}
	}

	return nil
}
