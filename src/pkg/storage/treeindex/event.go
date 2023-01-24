package treeindex

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"

	"github.com/pkg/errors"
)

func processHistoryItemStructure(node Noder, obj base.History[base.ItemStructure]) (Noder, error) {
	eventsNode, err := walk(obj.Events)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process HISTORY.Events item")
	}

	node.addAttribute("events", eventsNode)

	return node, nil
}

func processEventItemStructure(node Noder, obj base.Event[base.ItemStructure]) (Noder, error) {
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

func proccessPointEventItemStructure(node Noder, obj *base.PointEvent[base.ItemStructure]) (Noder, error) {
	node, err := proccessBaseEventItemStructure(node, &obj.BaseEvent)
	return node, err
}

func proccessIntervalEventItemStructure(node Noder, obj *base.IntervalEvent[base.ItemStructure]) (Noder, error) {
	node, err := proccessBaseEventItemStructure(node, &obj.BaseEvent)
	return node, errors.Wrap(err, "not implemented")
}

func proccessBaseEventItemStructure(node Noder, obj *base.BaseEvent[base.ItemStructure]) (Noder, error) {
	dataNode, err := walk(obj.Data)
	if err != nil {
		return nil, errors.Wrap(err, "cannot handle event data")
	}

	node.addAttribute("data", dataNode)

	return node, nil
}
