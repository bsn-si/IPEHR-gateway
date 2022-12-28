package treeindex

import (
	"encoding/json"
)

func (t *Tree) Print() string {
	data, _ := json.MarshalIndent(t.data, "", "    ")
	return string(data)
}
