package treeindex

import (
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
)

func processComposition(cmp model.Composition) (Noder, error) {
	node := newCompositionNode(cmp)

	node.addAttribute("language", newNode(&cmp.Language))
	node.addAttribute("territory", newNode(&cmp.Territory))
	node.addAttribute("category", newNode(&cmp.Category))

	if cmp.Context != nil {
		ctxNode, err := walk(*cmp.Context)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get node for Composition.Context")
		}

		node.addAttribute("context", ctxNode)
	}

	if err := node.processCompositionContent(cmp.Content); err != nil {
		return nil, errors.Wrap(err, "cannot process Composition.Content")
	}
	// TODO: add Composition.Composer field handling

	return node, nil
}

type CompositionNode struct {
	baseNode
	Tree
	attributes map[string]Noder
}

func newCompositionNode(cmp model.Composition) *CompositionNode {
	l := cmp.Locatable
	node := &CompositionNode{
		baseNode: baseNode{
			ID:   l.ArchetypeNodeID,
			Type: l.Type,
			Name: l.Name.Value,
		},
		Tree:       *NewTree(),
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
	n := cmp.baseNode.TryGetChild(key)
	if n != nil {
		return n
	}

	return cmp.attributes[key]
}

func (cmp CompositionNode) ForEach(foo func(name string, node Noder) bool) {
	for k, n := range cmp.attributes {
		for foo(k, n) {
		}
	}
}

func (cmp CompositionNode) addAttribute(key string, val Noder) {
	cmp.attributes[key] = val
}
