package treeindex

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (t *Tree) Print() string {
	builder := &strings.Builder{}

	printCollection(builder, "Actions", t.actions)
	printCollection(builder, "Evaluation", t.evaluations)
	printCollection(builder, "Instructions", t.instructions)
	printCollection(builder, "Observations", t.obeservations)

	return builder.String()
}

func printCollection(w io.Writer, name string, collection Container) {
	fmt.Fprintln(w, name)
	m := json.NewEncoder(w)
	m.SetIndent("", "\t")
	_ = m.Encode(collection)
	fmt.Fprintln(w)
}
