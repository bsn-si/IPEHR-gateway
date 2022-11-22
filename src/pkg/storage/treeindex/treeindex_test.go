package treeindex

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model"
	"os"
	"testing"
)

func TestTree_walk(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := wd + "/../../../../data/mock/ehr/composition.json"

	inJSON, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal("Can't open composition.json file", filePath)
	}

	c := model.Composition{}

	if err := json.Unmarshal(inJSON, &c); err != nil {
		t.Error(err)
		return
	}

	tree := NewTree()
	if err := tree.AddComposition(c); err != nil {
		t.Error(err)
	}

	// root, ok := tree.root..Children["openEHR-EHR-COMPOSITION.health_summary.v1"]
	// assert.Equal(t, ok, true)
	// assert.Equal(t, len(root.Children), 14)

	// t.Logf("%+v", tree.root.Children)
	t.Logf("tree:\n%s", tree.Print())
	t.Error("do error every time")
}
