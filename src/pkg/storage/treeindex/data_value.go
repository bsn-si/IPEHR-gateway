package treeindex

import (
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
)

func processDvValueBase(node noder, value *base.DvValueBase) (noder, error) {
	return node, nil
}

func processDvEncapsulated(node noder, value *base.DvEncapsulated) (noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_ENCAPSULATED.base")
	}

	if value.Charset != nil {
		node.addAttribute("charset", NewNodeForCodePhrase(*value.Charset))
	}

	if value.Language != nil {
		node.addAttribute("language", NewNodeForCodePhrase(*value.Language))
	}

	return node, nil
}

func processDvURI(node noder, value *base.DvURI) (noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_URI.base")
	}

	node.addAttribute("value", NewValueNode(value.Value))

	return node, nil
}

func processDvTemporal(node noder, value *base.DvTemporal) (noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_TEMPORAL.base")
	}

	if value.Accuracy != nil {
		newNode, err := processDvDuration(node, value.Accuracy)
		if err != nil {
			return nil, errors.Wrap(err, "cannot process DV_TEMPORAL.accuracy")
		}

		node.addAttribute("accuracy", newNode)
	}

	return node, nil
}

func processDvTime(node noder, value *base.DvTime) (noder, error) {
	node, err := processDvTemporal(node, &value.DvTemporal)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_TIME.base")
	}

	node.addAttribute("value", NewValueNode(value.Value))

	return node, nil
}

func processDvQuantity(node noder, value *base.DvQuantity) (noder, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_QUANTITY.base")
	}

	node.addAttribute("magnitude", NewValueNode(value.Magnitude))
	if value.Precision != nil {
		node.addAttribute("precision", NewValueNode(value.Precision))
	}

	if value.Units != nil {
		node.addAttribute("units", NewValueNode(value.Units))
	}

	if value.UnitsSystem != nil {
		node.addAttribute("units_system", NewValueNode(value.UnitsSystem))
	}

	if value.UnitsDisplayName != nil {
		node.addAttribute("units_display_name", NewValueNode(value.UnitsDisplayName))
	}

	//todo: add processing for DvQuantity.NormalRange and DvQuantity.OtherReferenceRanges fields

	return node, nil
}

func processDvState(node noder, value *base.DvState) (noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_STATE.base")
	}

	valueNode, err := walkDataValue(value.Value)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_STATE.value")
	}

	node.addAttribute("value", valueNode)
	node.addAttribute("is_terminal", NewValueNode(value.IsTerminal))

	return node, nil
}

func processDvOrdered[T any](node noder, value *base.DvOrdered[T]) (noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_ORDERED.base")
	}

	if value.NormalStatus != nil {
		node.addAttribute("norma_status", NewNodeForCodePhrase(*value.NormalStatus))
	}

	//todo: add processing for DvOrdered.NormalRange and DvOrdered.OtherReferenceRanges fields

	return node, nil
}

func processDvQuantified[T any](node noder, value *base.DvQuantified[T]) (noder, error) {
	node, err := processDvOrdered(node, &value.DvOrdered)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_QUANTIFIED.base")
	}

	node.addAttribute("magnitude_status", NewValueNode(value.MagnitudeStatus))

	//todo: add processing for value.Accuracy field

	return node, nil
}

func processDvAmount[T any](node noder, value *base.DvAmount[T]) (noder, error) {
	node, err := processDvQuantified(node, &value.DvQuantified)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_AMOUNT.base")
	}

	node.addAttribute("accuracy_is_percent", NewValueNode(value.AccuracyIsPercent))
	node.addAttribute("accuracy", NewValueNode(value.Accuracy))

	return node, nil
}

func processDvProportion(node noder, value *base.DvProportion) (noder, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_PROPORTION.base")
	}

	node.addAttribute("numeration", NewValueNode(value.Numeration))
	node.addAttribute("denomination", NewValueNode(value.Denomination))
	node.addAttribute("type", NewValueNode(value.Type))
	if value.Precision != nil {
		node.addAttribute("precision", NewValueNode(value.Precision))
	}

	//todo: add processing for value.NormalRange and value.OtherReferenceRanges

	return node, nil
}

func processDvParsable(node noder, value *base.DvParsable) (noder, error) {
	return node, errors.New("DV_PARSABLE not implemented")
}

func processDvParagraph(node noder, value *base.DvParagraph) (noder, error) {
	return node, errors.New("DV_PARAGRAPH not implemented")
}

func processDvMultimedia(node noder, value *base.DvMultimedia) (noder, error) {
	var err error

	node, err = processDvEncapsulated(node, &value.DvEncapsulated)
	if err != nil {
		return nil, err
	}

	if value.URI != nil {
		uriNode, err := walkDataValue(value.URI)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get node for DV_MULTIMEDIA.uri")
		}

		node.addAttribute("uri", uriNode)
	}

	node.addAttribute("alternate_text", NewValueNode(value.AlternativeText))
	node.addAttribute("size", NewValueNode(value.Size))
	node.addAttribute("media_type", NewNodeForCodePhrase(value.MediaType))
	if value.CompressionAlgorithm != nil {
		node.addAttribute("compression_algorithm", NewNodeForCodePhrase(*value.CompressionAlgorithm))
	}

	if value.IntegrityCheckAlgorithm != nil {
		node.addAttribute("integrity_check_algorithm", NewNodeForCodePhrase(*value.IntegrityCheckAlgorithm))
	}

	if value.IntegrityCheck != nil {
		node.addAttribute("integrity_check", NewValueNode(value.IntegrityCheck))
	}

	if value.Thumbnail != nil {
		thumbnailNode, err := walkDataValue(value.Thumbnail)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get node for DV_MULTIMEDIA.thumbnail")
		}

		node.addAttribute("thumbnail", thumbnailNode)
	}

	return node, nil
}

func processDvIdentifier(node noder, value *base.DvIdentifier) (noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_IDENTIFIER.base")
	}

	node.addAttribute("issuser", NewValueNode(value.Issuer))
	node.addAttribute("assigner", NewValueNode(value.Assigner))
	node.addAttribute("id", NewValueNode(value.ID))
	node.addAttribute("type", NewValueNode(value.Type))

	return node, nil
}

func processDvDuration(node noder, value *base.DvDuration) (noder, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_DURATION.base")
	}

	node.addAttribute("value", NewValueNode(value.Value))

	return node, nil
}

func processDvDateTime(node noder, value *base.DvDateTime) (noder, error) {
	node, err := processDvTemporal(node, &value.DvTemporal)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_DATE_TIME.base")
	}

	node.addAttribute("value", NewValueNode(value.Value))

	return node, nil
}

func processDvDate(node noder, value *base.DvDate) (noder, error) {
	node, err := processDvTemporal(node, &value.DvTemporal)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_DATE.base")
	}

	node.addAttribute("value", NewValueNode(value.Value))

	return node, nil
}

func processDvCount(node noder, value *base.DvCount) (noder, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_COUNT.base")
	}

	node.addAttribute("magnitude", NewValueNode(value.Magnitude))

	//todo: add processing for DV_COUNT.normal_range and DV_COUNT.other_reference_ranges fields

	return node, nil
}

func processDvCodedText(node noder, value *base.DvCodedText) (noder, error) {
	node, err := processDvText(node, &value.DvText)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_CODED_TEXT.DV_TEXT")
	}

	node.addAttribute("defining_code", NewNodeForCodePhrase(value.DefiningCode))

	return node, nil
}

func processDvText(node noder, value *base.DvText) (noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_TEXT.base")
	}

	node.addAttribute("value", NewValueNode(value.Value))
	if value.Formatting != "" {
		node.addAttribute("formatting", NewValueNode(value.Formatting))
	}

	if value.Hyperlink != nil {
		hyperlinkNode, err := processDvURI(node, value.Hyperlink)
		if err != nil {
			return nil, errors.Wrap(err, "cannot process DV_TEXT.hyperlink")
		}

		node.addAttribute("hyperlink", hyperlinkNode)
	}

	//todo: add processing for mappings filed

	if value.Language != nil {
		node.addAttribute("language", NewNodeForCodePhrase(*value.Language))
	}

	if value.Encoding != nil {
		node.addAttribute("encoding", NewNodeForCodePhrase(*value.Encoding))
	}

	return node, nil
}

func processDvBoolean(node noder, value *base.DvBoolean) (noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_BOLLEAN.base")
	}

	node.addAttribute("value", NewValueNode(value.Value))

	return node, nil
}
