package treeindex

import (
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
)

func processElement(node *Node, obj *base.Element) (*Node, error) {
	if obj.Value != nil {
		valueNode, err := walkDataValue(*obj.Value)
		if err != nil {
			return nil, errors.Wrap(err, "cannot handle ELEMENT.Value object")
		}

		node.Attributes.add("value", valueNode)
	}

	return node, nil
}

func processCluster(node *Node, obj *base.Cluster) (*Node, error) {
	for _, item := range obj.Items {
		itemNode, err := walk(item)
		if err != nil {
			return nil, errors.Wrap(err, "cannot handle CLUSTER.Items item object")
		}

		node.Attributes.add("items", itemNode)
	}

	return node, nil
}
