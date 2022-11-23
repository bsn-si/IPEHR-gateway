package treeindex

import (
	"encoding/json"
)

func (t *Tree) Print() string {
	m := map[string]any{
		"ACTIONS":      t.actions,
		"EVALUATIONS":  t.evaluations,
		"INSTRUCTIONS": t.instructions,
		"OBSERVATIONS": t.obeservations,
	}
	data, _ := json.MarshalIndent(m, "", "    ")

	return string(data)
}
