package treeindex

import (
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
)

func processElement(node noder, obj *base.Element) (noder, error) {
	if obj.Value != nil {
		valueNode, err := walkDataValue(*obj.Value)
		if err != nil {
			return nil, errors.Wrap(err, "cannot handle ELEMENT.Value object")
		}

		node.addAttribute("value", valueNode)
	}

	return node, nil
}

func processCluster(node noder, obj *base.Cluster) (noder, error) {
	for _, item := range obj.Items {
		itemNode, err := walk(item)
		if err != nil {
			return nil, errors.Wrap(err, "cannot handle CLUSTER.Items item object")
		}

		node.addAttribute("items", itemNode)
	}

	return node, nil
}
