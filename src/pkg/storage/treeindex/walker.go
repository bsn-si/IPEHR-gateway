package treeindex

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func walk(obj any) (Noder, error) {
	switch obj := obj.(type) {
	case *model.EHR:
		return ProcessEHR(obj)
	case *model.Composition:
		return ProcessComposition(obj)
	case model.EventContext:
		return processEventContext(obj)
	case base.Root:
		return walkRoot(obj)
	case base.DataValue:
		return walkDataValue(obj)
	default:
		return walkBySlice(obj)
	}
}

func walkRoot(obj base.Root) (Noder, error) {
	var err error
	node := newNode(obj) //nolint

	switch obj := obj.(type) {
	case *base.Action:
		node, err = processAction(node, obj)
	case *base.Evaluation:
		node, err = processEvaluation(node, obj)
	case *base.Instruction:
		node, err = processInstruction(node, obj)
	case *base.Observation:
		node, err = processObservation(node, obj)
	case base.History[base.ItemStructure]:
		node, err = processHistoryItemStructure(node, obj)
	case base.Event[base.ItemStructure]:
		node, err = processEventItemStructure(node, obj)
	case base.ItemStructure:
		node, err = processItemStructure(node, obj)
	case *base.Element:
		node, err = processElement(node, obj)
	case *base.Cluster:
		node, err = processCluster(node, obj)
	case base.ItemTree:
		node, err = processItemTree(node, &obj)
	default:
		return nil, fmt.Errorf("unexpected node type: %T", obj) //nolint
	}

	if err != nil {
		return nil, err
	}

	return node, nil
}

func walkBySlice(slice any) (Noder, error) {
	sliceNode := newSliceNode()

	switch ss := slice.(type) {
	case []model.Composition:
		for _, item := range ss {
			node, err := walk(item)
			if err != nil {
				return nil, errors.Wrap(err, "cannot process slice COMPOSITION")
			}

			sliceNode.addAttribute(node.GetID(), node)
		}
	case []base.Root:
		for _, item := range ss {
			node, err := walk(item)
			if err != nil {
				return nil, errors.Wrap(err, "cannot process slice ITEM")
			}

			sliceNode.addAttribute(node.GetID(), node)
		}
	case base.Items:
		for _, item := range ss {
			node, err := walk(item)
			if err != nil {
				return nil, errors.Wrap(err, "cannot process ITEMS")
			}

			sliceNode.addAttribute(node.GetID(), node)
		}
	case []base.Event[base.ItemStructure]:
		for _, item := range ss {
			node, err := walk(item)
			if err != nil {
				return nil, errors.Wrap(err, "cannot process EVENTS slice")
			}

			sliceNode.addAttribute(node.GetID(), node)
		}
	default:
		return nil, fmt.Errorf("unexpected slice type: %T", slice) //nolint
	}

	return sliceNode, nil
}

func walkDataValue(dv base.DataValue) (Noder, error) {
	var err error
	node := newNode(dv) //nolint

	switch value := dv.(type) {
	case *base.DvURI:
		node, err = processDvURI(node, value)
	case *base.DvTime:
		node, err = processDvTime(node, value)
	case *base.DvQuantity:
		node, err = processDvQuantity(node, value)
	case *base.DvState:
		node, err = processDvState(node, value)
	case *base.DvProportion:
		node, err = processDvProportion(node, value)
	case *base.DvParsable:
		node, err = processDvParsable(node, value)
	case *base.DvParagraph:
		node, err = processDvParagraph(node, value)
	case *base.DvMultimedia:
		node, err = processDvMultimedia(node, value)
	case *base.DvIdentifier:
		node, err = processDvIdentifier(node, value)
	case *base.DvDuration:
		node, err = processDvDuration(node, value)
	case *base.DvDateTime:
		node, err = processDvDateTime(node, value)
	case *base.DvDate:
		node, err = processDvDate(node, value)
	case *base.DvCount:
		node, err = processDvCount(node, value)
	case *base.DvCodedText:
		node, err = processDvCodedText(node, value)
	case *base.DvText:
		node, err = processDvText(node, value)
	case *base.DvBoolean:
		node, err = processDvBoolean(node, value)
	default:
		return nil, fmt.Errorf("unexpected value type: %T", value) //nolint
	}

	return node, err
}
