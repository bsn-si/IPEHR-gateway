package treeindex

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func ProcessComposition(cmp *model.Composition) (*CompositionNode, error) {
	node := newCompositionNode(cmp)

	node.addAttribute("language", newNode(cmp.Language))
	node.addAttribute("territory", newNode(cmp.Territory))
	node.addAttribute("category", newNode(cmp.Category))

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
	BaseNode
	Tree
	Attributes Attributes
}

func newCompositionNode(cmp *model.Composition) *CompositionNode {
	l := cmp.Locatable
	node := &CompositionNode{
		BaseNode: BaseNode{
			ID:       l.ArchetypeNodeID,
			Type:     l.Type,
			Name:     l.Name.Value,
			NodeType: CompostionNodeType,
		},
		Tree:       *NewTree(),
		Attributes: Attributes{},
	}

	return node
}

func (cmp CompositionNode) GetID() string {
	return cmp.ID
}

func (cmp CompositionNode) TryGetChild(key string) Noder {
	n := cmp.BaseNode.TryGetChild(key)
	if n != nil {
		return n
	}

	return cmp.Attributes[key]
}

func (cmp CompositionNode) addAttribute(key string, val Noder) {
	cmp.Attributes[key] = val
}
