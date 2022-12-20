package treeindex

import "hms/gateway/pkg/docs/model"

type CompositionNode struct {
}

func newCompositionNode(cmp model.Composition) Noder {
	return &CompositionNode{}
}

func (cmp CompositionNode) GetNodeType() NodeType {
	return CompostionNodeType
}

func (cmp CompositionNode) GetID() string {
	return ""
}

func (cmp CompositionNode) TryGetChild(key string) Noder {
	return nil
}

func (cmp CompositionNode) ForEach(func(name string, node Noder) bool) {
}

func (cmp CompositionNode) addAttribute(key string, val Noder) {

}
