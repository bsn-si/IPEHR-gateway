package treeindex

import (
	"fmt"
	"hms/gateway/pkg/docs/model/base"

	"github.com/pkg/errors"
)

func processHistoryItemStructure(node *Node, obj base.History[base.ItemStructure]) (*Node, error) {
	for _, e := range obj.Events {
		eventsNode, err := walk(e)
		if err != nil {
			return nil, errors.Wrap(err, "cannot process HISTORY.Events item")
		}
		node.Attributes.add("events", eventsNode)
	}
	return node, nil
}

func processEventItemStructure(node *Node, obj base.Event[base.ItemStructure]) (*Node, error) {
	switch obj.GetType() {
	case base.PointEventItemType:
		pointEvent, ok := obj.Data.(*base.PointEvent[base.ItemStructure])
		if !ok {
			return nil, fmt.Errorf("Event.Data invalid type: %T", obj.Data) // nolint
		}

		return proccessPointEventItemStructure(node, pointEvent)
	case base.IntervalEventItemType:
		intervalEvent, ok := obj.Data.(*base.IntervalEvent[base.ItemStructure])
		if !ok {
			return nil, fmt.Errorf("Event.Data invalid type: %T", obj.Data) // nolint
		}

		return proccessIntervalEventItemStructure(node, intervalEvent)

	default:
		return nil, fmt.Errorf("unexpected event type: %v", obj.GetType()) // nolint
	}
}

func proccessPointEventItemStructure(node *Node, obj *base.PointEvent[base.ItemStructure]) (*Node, error) {
	node, err := proccessBaseEventItemStructure(node, &obj.BaseEvent)
	return node, err
}

func proccessIntervalEventItemStructure(node *Node, obj *base.IntervalEvent[base.ItemStructure]) (*Node, error) {
	node, err := proccessBaseEventItemStructure(node, &obj.BaseEvent)
	return node, errors.Wrap(err, "not implemented")
}

func proccessBaseEventItemStructure(node *Node, obj *base.BaseEvent[base.ItemStructure]) (*Node, error) {
	dataNode, err := walk(obj.Data)
	if err != nil {
		return nil, errors.Wrap(err, "cannot handle event data")
	}

	node.Attributes.add("data", dataNode)

	return node, nil
}
