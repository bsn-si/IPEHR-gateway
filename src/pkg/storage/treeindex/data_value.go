package treeindex

import (
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
)

func processDvValueBase(node *Node, value *base.DvValueBase) (*Node, error) {
	return node, nil
}

func processDvEncapsulated(node *Node, value *base.DvEncapsulated) (*Node, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_ENCAPSULATED.base")
	}

	if value.Charset != nil {
		node.Value["charset"] = NewNodeForCodePhrase(*value.Charset)
	}

	if value.Language != nil {
		node.Value["language"] = NewNodeForCodePhrase(*value.Language)
	}

	return node, nil
}

func processDvURI(node *Node, value *base.DvURI) (*Node, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_URI.base")
	}

	node.Value["value"] = value.Value

	return node, nil
}

func processDvTemporal(node *Node, value *base.DvTemporal) (*Node, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_TEMPORAL.base")
	}

	if value.Accuracy != nil {
		node.Value["accuracy"], err = processDvDuration(node, value.Accuracy)
		if err != nil {
			return nil, errors.Wrap(err, "cannot process DV_TEMPORAL.accuracy")
		}
	}

	return node, nil
}

func processDvTime(node *Node, value *base.DvTime) (*Node, error) {
	node, err := processDvTemporal(node, &value.DvTemporal)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_TIME.base")
	}

	node.Value["value"] = value.Value

	return node, nil
}

func processDvQuantity(node *Node, value *base.DvQuantity) (*Node, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_QUANTITY.base")
	}

	node.Value["magnitude"] = value.Magnitude
	if value.Precision != nil {
		node.Value["precision"] = value.Precision
	}

	if value.Units != nil {
		node.Value["units"] = value.Units
	}

	if value.UnitsSystem != nil {
		node.Value["units_system"] = value.UnitsSystem
	}

	if value.UnitsDisplayName != nil {
		node.Value["units_display_name"] = value.UnitsDisplayName
	}

	//todo: add processing for DvQuantity.NormalRange and DvQuantity.OtherReferenceRanges fields

	return node, nil
}

func processDvState(node *Node, value *base.DvState) (*Node, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_STATE.base")
	}

	node.Value["value"], err = walkDataValue(value.Value)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_STATE.value")
	}

	node.Value["is_terminal"] = value.IsTerminal

	return node, nil
}

func processDvOrdered[T any](node *Node, value *base.DvOrdered[T]) (*Node, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_ORDERED.base")
	}

	if value.NormalStatus != nil {
		node.Value["norma_status"] = NewNodeForCodePhrase(*value.NormalStatus)
	}

	//todo: add processing for DvOrdered.NormalRange and DvOrdered.OtherReferenceRanges fields

	return node, nil
}

func processDvQuantified[T any](node *Node, value *base.DvQuantified[T]) (*Node, error) {
	node, err := processDvOrdered(node, &value.DvOrdered)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_QUANTIFIED.base")
	}

	node.Value["magnitude_status"] = value.MagnitudeStatus

	//todo: add processing for value.Accuracy field
	return node, nil
}

func processDvAmount[T any](node *Node, value *base.DvAmount[T]) (*Node, error) {
	node, err := processDvQuantified(node, &value.DvQuantified)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_AMOUNT.base")
	}

	node.Value["accuracy_is_percent"] = value.AccuracyIsPercent
	node.Value["accuracy"] = value.Accuracy

	return node, nil
}

func processDvProportion(node *Node, value *base.DvProportion) (*Node, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_PROPORTION.base")
	}

	node.Value["numeration"] = value.Numeration
	node.Value["denomination"] = value.Denomination
	node.Value["type"] = value.Type
	if value.Precision != nil {
		node.Value["precision"] = value.Precision
	}

	//todo: add processing for value.NormalRange and value.OtherReferenceRanges

	return node, nil
}

func processDvParsable(node *Node, value *base.DvParsable) (*Node, error) {
	return node, errors.New("DV_PARSABLE not implemented")
}

func processDvParagraph(node *Node, value *base.DvParagraph) (*Node, error) {
	return node, errors.New("DV_PARAGRAPH not implemented")
}

func processDvMultimedia(node *Node, value *base.DvMultimedia) (*Node, error) {
	var err error

	node, err = processDvEncapsulated(node, &value.DvEncapsulated)
	if err != nil {
		return nil, err
	}

	if value.URI != nil {
		node.Value["uri"], err = walkDataValue(value.URI)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get node for DV_MULTIMEDIA.uri")
		}
	}

	node.Value["alternate_text"] = value.AlternativeText
	node.Value["size"] = value.Size
	node.Value["media_type"] = NewNodeForCodePhrase(value.MediaType)
	if value.CompressionAlgorithm != nil {
		node.Value["compression_algorithm"] = NewNodeForCodePhrase(*value.CompressionAlgorithm)
	}

	if value.IntegrityCheckAlgorithm != nil {
		node.Value["integrity_check_algorithm"] = NewNodeForCodePhrase(*value.IntegrityCheckAlgorithm)
	}

	if value.IntegrityCheck != nil {
		node.Value["integrity_check"] = value.IntegrityCheck
	}

	if value.Thumbnail != nil {
		node.Value["thumbnail"], err = walkDataValue(value.Thumbnail)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get node for DV_MULTIMEDIA.thumbnail")
		}
	}
	return node, nil
}

func processDvIdentifier(node *Node, value *base.DvIdentifier) (*Node, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_IDENTIFIER.base")
	}

	node.Value["issuser"] = value.Issuer
	node.Value["assigner"] = value.Assigner
	node.Value["id"] = value.ID
	node.Value["type"] = value.Type

	return node, nil
}

func processDvDuration(node *Node, value *base.DvDuration) (*Node, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_DURATION.base")
	}

	node.Value["value"] = value.Value

	return node, nil
}

func processDvDateTime(node *Node, value *base.DvDateTime) (*Node, error) {
	node, err := processDvTemporal(node, &value.DvTemporal)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_DATE_TIME.base")
	}

	node.Value["value"] = value.Value

	return node, nil
}

func processDvDate(node *Node, value *base.DvDate) (*Node, error) {
	node, err := processDvTemporal(node, &value.DvTemporal)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_DATE.base")
	}

	node.Value["value"] = value.Value

	return node, nil
}

func processDvCount(node *Node, value *base.DvCount) (*Node, error) {
	node, err := processDvAmount(node, &value.DvAmount)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_COUNT.base")
	}

	node.Value["magnitude"] = value.Magnitude

	//todo: add processing for DV_COUNT.normal_range and DV_COUNT.other_reference_ranges fields

	return node, nil
}

func processDvCodedText(node *Node, value *base.DvCodedText) (*Node, error) {
	node, err := processDvText(node, &value.DvText)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_CODED_TEXT.DV_TEXT")
	}

	node.Value["defining_code"] = NewNodeForCodePhrase(value.DefiningCode)

	return node, nil
}

func processDvText(node *Node, value *base.DvText) (*Node, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_TEXT.base")
	}

	node.Value["value"] = value.Value
	if value.Formatting != "" {
		node.Value["formatting"] = value.Formatting
	}

	if value.Hyperlink != nil {
		node.Value["hyperlink"], err = processDvURI(node, value.Hyperlink)
		if err != nil {
			return nil, errors.Wrap(err, "cannot process DV_TEXT.hyperlink")
		}
	}

	//todo: add processing for mappings filed

	if value.Language != nil {
		node.Value["language"] = NewNodeForCodePhrase(*value.Language)
	}

	if value.Encoding != nil {
		node.Value["encoding"] = NewNodeForCodePhrase(*value.Encoding)
	}

	return node, nil
}

func processDvBoolean(node *Node, value *base.DvBoolean) (*Node, error) {
	node, err := processDvValueBase(node, &value.DvValueBase)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process DV_BOLLEAN.base")
	}

	node.Value["value"] = value.Value
	return node, nil
}
