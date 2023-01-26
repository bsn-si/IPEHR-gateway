package treeindex

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func processDvValueBase(node Noder, value *base.DvValueBase) (Noder, error) {
	return node, nil
}

func processDvEncapsulated(node Noder, value *base.DvEncapsulated) (Noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_ENCAPSULATED.base")
	}

	if value.Charset != nil {
		node.addAttribute("charset", newNode(*value.Charset))
	}

	if value.Language != nil {
		node.addAttribute("language", newNode(*value.Language))
	}

	return node, nil
}

func processDvURI(node Noder, value *base.DvURI) (Noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_URI.base")
	}

	node.addAttribute("value", newNode(value.Value))

	return node, nil
}

func processDvTemporal(node Noder, value *base.DvTemporal) (Noder, error) {
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

func processDvTime(node Noder, value *base.DvTime) (Noder, error) {
	node, err := processDvTemporal(node, &value.DvTemporal)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_TIME.base")
	}

	node.addAttribute("value", newNode(value.Value))

	return node, nil
}

func processDvQuantity(node Noder, value *base.DvQuantity) (Noder, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_QUANTITY.base")
	}

	node.addAttribute("magnitude", newNode(value.Magnitude))

	if value.Precision != nil {
		node.addAttribute("precision", newNode(*value.Precision))
	}

	if value.Units != nil {
		node.addAttribute("units", newNode(*value.Units))
	}

	if value.UnitsSystem != nil {
		node.addAttribute("units_system", newNode(*value.UnitsSystem))
	}

	if value.UnitsDisplayName != nil {
		node.addAttribute("units_display_name", newNode(*value.UnitsDisplayName))
	}

	//todo: add processing for DvQuantity.NormalRange and DvQuantity.OtherReferenceRanges fields

	return node, nil
}

func processDvState(node Noder, value *base.DvState) (Noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_STATE.base")
	}

	valueNode, err := walk(value.Value)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_STATE.value")
	}

	node.addAttribute("value", valueNode)
	node.addAttribute("is_terminal", newNode(value.IsTerminal))

	return node, nil
}

func processDvOrdered[T any](node Noder, value *base.DvOrdered[T]) (Noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_ORDERED.base")
	}

	if value.NormalStatus != nil {
		node.addAttribute("norma_status", newNode(*value.NormalStatus))
	}

	//todo: add processing for DvOrdered.NormalRange and DvOrdered.OtherReferenceRanges fields

	return node, nil
}

func processDvQuantified[T any](node Noder, value *base.DvQuantified[T]) (Noder, error) {
	node, err := processDvOrdered(node, &value.DvOrdered)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_QUANTIFIED.base")
	}

	node.addAttribute("magnitude_status", newNode(value.MagnitudeStatus))

	//todo: add processing for value.Accuracy field

	return node, nil
}

func processDvAmount[T any](node Noder, value *base.DvAmount[T]) (Noder, error) {
	node, err := processDvQuantified(node, &value.DvQuantified)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_AMOUNT.base")
	}

	node.addAttribute("accuracy_is_percent", newNode(value.AccuracyIsPercent))
	node.addAttribute("accuracy", newNode(value.Accuracy))

	return node, nil
}

func processDvProportion(node Noder, value *base.DvProportion) (Noder, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_PROPORTION.base")
	}

	node.addAttribute("numeration", newNode(value.Numeration))
	node.addAttribute("denomination", newNode(value.Denomination))
	node.addAttribute("type", newNode(value.Type))

	if value.Precision != nil {
		node.addAttribute("precision", newNode(*value.Precision))
	}

	//todo: add processing for value.NormalRange and value.OtherReferenceRanges

	return node, nil
}

func processDvParsable(node Noder, value *base.DvParsable) (Noder, error) {
	return node, errors.New("DV_PARSABLE not implemented")
}

func processDvParagraph(node Noder, value *base.DvParagraph) (Noder, error) {
	return node, errors.New("DV_PARAGRAPH not implemented")
}

func processDvMultimedia(node Noder, value *base.DvMultimedia) (Noder, error) {
	var err error

	node, err = processDvEncapsulated(node, &value.DvEncapsulated)
	if err != nil {
		return nil, err
	}

	if value.URI != nil {
		uriNode, err := walk(value.URI)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get node for DV_MULTIMEDIA.uri")
		}

		node.addAttribute("uri", uriNode)
	}

	node.addAttribute("alternate_text", newNode(value.AlternativeText))
	node.addAttribute("size", newNode(value.Size))
	node.addAttribute("media_type", newNode(value.MediaType))

	if value.CompressionAlgorithm != nil {
		node.addAttribute("compression_algorithm", newNode(*value.CompressionAlgorithm))
	}

	if value.IntegrityCheckAlgorithm != nil {
		node.addAttribute("integrity_check_algorithm", newNode(*value.IntegrityCheckAlgorithm))
	}

	if value.IntegrityCheck != nil {
		node.addAttribute("integrity_check", newNode(value.IntegrityCheck))
	}

	if value.Thumbnail != nil {
		thumbnailNode, err := walk(value.Thumbnail)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get node for DV_MULTIMEDIA.thumbnail")
		}

		node.addAttribute("thumbnail", thumbnailNode)
	}

	return node, nil
}

func processDvIdentifier(node Noder, value *base.DvIdentifier) (Noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_IDENTIFIER.base")
	}

	node.addAttribute("issuser", newNode(value.Issuer))
	node.addAttribute("assigner", newNode(value.Assigner))
	node.addAttribute("id", newNode(value.ID))
	node.addAttribute("type", newNode(value.Type))

	return node, nil
}

func processDvDuration(node Noder, value *base.DvDuration) (Noder, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_DURATION.base")
	}

	node.addAttribute("value", newNode(value.Value))

	return node, nil
}

func processDvDateTime(node Noder, value *base.DvDateTime) (Noder, error) {
	node, err := processDvTemporal(node, &value.DvTemporal)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_DATE_TIME.base")
	}

	node.addAttribute("value", newNode(value.Value))

	return node, nil
}

func processDvDate(node Noder, value *base.DvDate) (Noder, error) {
	node, err := processDvTemporal(node, &value.DvTemporal)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_DATE.base")
	}

	node.addAttribute("value", newNode(value.Value))

	return node, nil
}

func processDvCount(node Noder, value *base.DvCount) (Noder, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_COUNT.base")
	}

	node.addAttribute("magnitude", newNode(value.Magnitude))

	//todo: add processing for DV_COUNT.normal_range and DV_COUNT.other_reference_ranges fields

	return node, nil
}

func processDvCodedText(node Noder, value *base.DvCodedText) (Noder, error) {
	node, err := processDvText(node, &value.DvText)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_CODED_TEXT.DV_TEXT")
	}

	node.addAttribute("defining_code", newNode(value.DefiningCode))

	return node, nil
}

func processDvText(node Noder, value *base.DvText) (Noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_TEXT.base")
	}

	node.addAttribute("value", newNode(value.Value))
	node.addAttribute("formatting", newNode(value.Formatting))

	if value.Hyperlink != nil {
		hyperlinkNode, err := processDvURI(node, value.Hyperlink)
		if err != nil {
			return nil, errors.Wrap(err, "cannot process DV_TEXT.hyperlink")
		}

		node.addAttribute("hyperlink", hyperlinkNode)
	}

	//todo: add processing for mappings filed

	if value.Language != nil {
		node.addAttribute("language", newNode(*value.Language))
	}

	if value.Encoding != nil {
		node.addAttribute("encoding", newNode(*value.Encoding))
	}

	return node, nil
}

func processDvBoolean(node Noder, value *base.DvBoolean) (Noder, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_BOLLEAN.base")
	}

	node.addAttribute("value", newNode(value.Value))

	return node, nil
}
