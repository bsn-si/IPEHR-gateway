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
	case *base.Observation:
		node, err = processObservation(node, obj)
	case base.History[base.ItemStructure]:
		node, err = processHistoryItemStructure(node, obj)
	case base.Event[base.ItemStructure]:
		node, err = processEventItemStructure(node, obj)
	case base.ItemStructure:
		node, err = processItemStructure(node, obj)
	default:
		return nil, fmt.Errorf("unexpected node type: %T", obj) //nolint
	}

	if err != nil {
		return nil, err
	}

	return node, nil
}

func processObservation(node *Node, obs *base.Observation) (*Node, error) {
	dataNode, err := walk(obs.Data)
	if err != nil {
		return nil, errors.Wrap(err, "cannon process OBSERVATION.Data")
	}
	node.Attributes.add("data", dataNode)

	if obs.Protocol.Data != nil {
		protocolNode, err := walk(obs.Protocol)
		if err != nil {
			return nil, errors.Wrap(err, "cannon process OBSERVATION.Protocol")
		}
		node.Attributes.add("protocol", protocolNode)
	}

	return node, nil
}
