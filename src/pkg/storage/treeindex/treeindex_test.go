package treeindex

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
	err = tree.AddComposition(c)
	assert.Nil(t, err)

	assert.Equal(t, 4, tree.actions.Len())
	assert.Equal(t, 21, tree.evaluations.Len())
	assert.Equal(t, 1, tree.instructions.Len())
	assert.Equal(t, 12, tree.obeservations.Len())

	treeStr := tree.Print()
	// uncomment next string tp view tre string
	// t.Error(treeStr)
	if treeStr == "" {
		t.Fail()
	}
}
