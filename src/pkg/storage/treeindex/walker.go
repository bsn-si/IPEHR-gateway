package treeindex

import (
	"fmt"
	"hms/gateway/pkg/docs/model/base"

	"github.com/pkg/errors"
)

func walk(obj base.Root) (*Node, error) {
	var err error
	node := NewNode(obj)

	switch obj := obj.(type) {
	case *base.Action:
		node, err = processAction(node, obj)
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

func walkDataValue(dv base.DataValue) (*Node, error) {
	var err error
	node := NewNodeForData(dv)

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

func processAction(node *Node, act *base.Action) (*Node, error) {
	descriptionNode, err := walk(act.Description)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ACTION.description")
	}
	node.Attributes.add("description", descriptionNode)

	return node, nil
}

func processObservation(node *Node, obs *base.Observation) (*Node, error) {
	dataNode, err := walk(obs.Data)
	if err != nil {
		return nil, errors.Wrap(err, "cannon process OBSERVATION.Data")
	}
	node.Attributes.add("data", dataNode)

	stateNode, err := walk(obs.State)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process OBSERVATION.state")
	}
	node.Attributes.add("state", stateNode)

	if obs.Protocol.Data != nil {
		protocolNode, err := walk(obs.Protocol)
		if err != nil {
			return nil, errors.Wrap(err, "cannon process OBSERVATION.Protocol")
		}
		node.Attributes.add("protocol", protocolNode)
	}

	return node, nil
}
