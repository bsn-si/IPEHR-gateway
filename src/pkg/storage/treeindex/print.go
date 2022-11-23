package treeindex

import (
	"encoding/json"
)

func (t *Tree) Print() string {
	result := "Observations:\n"
	obs, _ := json.MarshalIndent(t.obeservations, "", "  ")
	result += string(obs)
	return result
}
