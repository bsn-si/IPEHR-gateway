package treeindex

import (
	"fmt"
	"hms/gateway/pkg/docs/model/base"

	"github.com/pkg/errors"
)

func processHistoryItemStructure(node noder, obj base.History[base.ItemStructure]) (noder, error) {
	for _, e := range obj.Events {
		eventsNode, err := walk(e)
		if err != nil {
			return nil, errors.Wrap(err, "cannot process HISTORY.Events item")
		}

		node.addAttribute("events", eventsNode)
	}

	return node, nil
}

func processEventItemStructure(node noder, obj base.Event[base.ItemStructure]) (noder, error) {
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

func proccessPointEventItemStructure(node noder, obj *base.PointEvent[base.ItemStructure]) (noder, error) {
	node, err := proccessBaseEventItemStructure(node, &obj.BaseEvent)
	return node, err
}

func proccessIntervalEventItemStructure(node noder, obj *base.IntervalEvent[base.ItemStructure]) (noder, error) {
	node, err := proccessBaseEventItemStructure(node, &obj.BaseEvent)
	return node, errors.Wrap(err, "not implemented")
}

func proccessBaseEventItemStructure(node noder, obj *base.BaseEvent[base.ItemStructure]) (noder, error) {
	dataNode, err := walk(obj.Data)
	if err != nil {
		return nil, errors.Wrap(err, "cannot handle event data")
	}

	node.addAttribute("data", dataNode)

	return node, nil
}
