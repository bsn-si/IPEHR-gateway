package treeindex

import (
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
)

func processElement(node noder, obj *base.Element) (noder, error) {
	if obj.Value != nil {
		valueNode, err := walk(*obj.Value)
		if err != nil {
			return nil, errors.Wrap(err, "cannot handle ELEMENT.Value object")
		}

		node.addAttribute("value", valueNode)
	}

	return node, nil
}

func processCluster(node noder, obj *base.Cluster) (noder, error) {
	itemsNode, err := walk(obj.Items)
	if err != nil {
		return nil, errors.Wrap(err, "cannot handle CLUSTER.Items")
	}

	node.addAttribute("items", itemsNode)
	return node, nil
}
