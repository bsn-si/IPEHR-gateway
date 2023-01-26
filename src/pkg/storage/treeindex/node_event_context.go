package treeindex

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
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
	BaseNode
	Attributes Attributes `json:"-"`
}

func NewEventContextNode(ctx model.EventContext) *EventContextNode {
	node := EventContextNode{
		BaseNode: BaseNode{
			NodeType: EventContextNodeType,
		},
		Attributes: Attributes{},
	}

	return &node
}

func (n *EventContextNode) GetID() string {
	return ""
}

func (n *EventContextNode) addAttribute(key string, val Noder) {
	n.Attributes[key] = val
}

func (n *EventContextNode) TryGetChild(key string) Noder {
	return n.Attributes[key]
}
