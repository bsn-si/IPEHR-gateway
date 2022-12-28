package treeindex

import (
	"hms/gateway/pkg/docs/model"
)

func processEventContext(ctx model.EventContext) (Noder, error) {
	node := NewEventContextNode(ctx)

	node.addAttribute("start_time", newNode(&ctx.StartTime))

	if ctx.EndTime != nil {
		node.addAttribute("end_time", newNode(ctx.EndTime))
	}

	if ctx.Location != nil {
		node.addAttribute("location", newNode(*ctx.Location))
	}

	node.addAttribute("setting", newNode(ctx.Setting))

	// TODO: add OtherContext handling
	// if ctx.OtherContext != nil {
	// 	node.addAttribute("other_context", newNode(ctx.OtherContext))
	// }

	// TODO: add HealthCareFacility handling
	// if ctx.HealthCareFacility != nil {
	// 	hcfNode, err := walk(ctx.HealthCareFacility)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "cannot get EventContext.HealtCareFacility node")
	// 	}
	//
	// 	node.addAttribute("health_care_facility", hcfNode)
	// }

	// TODO: add Participations handling
	// if len(ctx.Participations) != 0 {
	// 	pNode, err := walk(ctx.Participations)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "cannot get EntentContext.Partisipations node")
	// 	}

	// 	node.addAttribute("participatios", pNode)
	// }

	return node, nil
}

type EventContextNode struct {
	attributes map[string]Noder `json:"-"`
}

func NewEventContextNode(ctx model.EventContext) *EventContextNode {
	node := EventContextNode{
		attributes: map[string]Noder{},
	}

	return &node
}

func (n *EventContextNode) GetNodeType() NodeType {
	return EventContextNodeType
}

func (n *EventContextNode) GetID() string {
	return ""
}

func (n *EventContextNode) addAttribute(key string, val Noder) {
	n.attributes[key] = val
}

func (n *EventContextNode) TryGetChild(key string) Noder {
	return n.attributes[key]
}

func (n *EventContextNode) ForEach(f func(name string, node Noder) bool) {
	for k, v := range n.attributes {
		if !f(k, v) {
			break
		}
	}
}
