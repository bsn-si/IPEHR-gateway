package treeindex

import (
	"fmt"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
)

type Tree struct {
	root *Node
}

type Node struct {
	Value    base.Root
	Children map[string]*Node
}

func NewTree() *Tree {
	return &Tree{
		root: &Node{},
	}
}

func (t *Tree) ContainsComposition(id string) bool {
	return false
}

func (t *Tree) AddComposition(cmp base.Root) {
	t.walk(cmp)
}

func (t *Tree) walk(obj base.Root) {
	switch t := obj.(type) {
	case model.Composition:
		{
			fmt.Println(t.ArchetypeNodeID)
		}
	default:
		fmt.Println("unexpected node")
	}
}
