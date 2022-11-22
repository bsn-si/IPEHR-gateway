package treeindex

import (
	"fmt"
	"strings"
)

func (t *Tree) Print() string {
	result := "Observations:\n"
	for key, container := range t.obeservations {
		result += fmt.Sprintf("\t[%s]\n", key)
		result += container.print(2)
	}

	return result
}

func (c Container) print(offset int) string {
	builder := &strings.Builder{}
	for _, nodes := range c {
		for _, node := range nodes {
			fmt.Fprint(builder, node.print(offset+1))
		}
	}

	return builder.String()
}

func (node *Node) print(offset int) string {
	offsetStr := strings.Repeat("\t", offset)

	builder := &strings.Builder{}
	fmt.Fprintf(builder, "%s[%s]-[%v]\n", offsetStr, node.Type, node.ID)

	if len(node.Attributes) > 0 {
		fmt.Fprint(builder, node.Attributes.print(offset+1))
	}

	if len(node.Value) > 0 {
		offsetStr := strings.Repeat("\t", offset+1)
		for key, val := range node.Value {
			fmt.Fprintf(builder, "%s%s: %v\n", offsetStr, key, val)
		}
	}

	return builder.String()
}

func (a Attributes) print(offset int) string {
	builder := &strings.Builder{}

	offsetStr := strings.Repeat("\t", offset)

	for attrName, attr := range a {
		fmt.Fprintf(builder, "%s%s\n", offsetStr, attrName)
		for _, node := range attr {
			fmt.Fprint(builder, node.print(offset+1))
		}
	}

	return builder.String()
}
