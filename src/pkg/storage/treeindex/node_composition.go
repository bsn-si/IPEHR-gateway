package treeindex

import "hms/gateway/pkg/docs/model"

func processComposition(cmp model.Composition) (Noder, error) {
	node := newCompositionNode(cmp)
	return node, nil
}

type CompositionNode struct {
	baseNode

	attributes map[string]Noder
}

func newCompositionNode(cmp model.Composition) Noder {
	l := cmp.Locatable
	node := &CompositionNode{
		baseNode: baseNode{
			ID:   l.ArchetypeNodeID,
			Type: l.Type,
			Name: l.Name.Value,
		},
		attributes: map[string]Noder{},
	}

	return node
}

func (cmp CompositionNode) GetNodeType() NodeType {
	return CompostionNodeType
}

func (cmp CompositionNode) GetID() string {
	return cmp.ID
}

func (cmp CompositionNode) TryGetChild(key string) Noder {
	return nil
}

func (cmp CompositionNode) ForEach(foo func(name string, node Noder) bool) {
	for k, n := range cmp.attributes {
		for foo(k, n) {
		}
	}
}

func (cmp CompositionNode) addAttribute(key string, val Noder) {

}
