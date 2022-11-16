package treeindex

import (
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"testing"
)

func TestTree_walk(t *testing.T) {
	c := model.Composition{
		Locatable: base.Locatable{
			ArchetypeNodeID: "some_composition_node_id",
		},
	}

	tree := NewTree()
	tree.AddComposition(c)

	if !tree.ContainsComposition(c.ArchetypeNodeID) {
		t.Fail()
	}
}
