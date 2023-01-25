package treeindex

import (
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func processElement(node Noder, obj *base.Element) (Noder, error) {
	if obj.Value != nil {
		valueNode, err := walk(obj.Value)
		if err != nil {
			return nil, errors.Wrap(err, "cannot handle ELEMENT.Value object")
		}

		node.addAttribute("value", valueNode)
	}

	return node, nil
}

func processCluster(node Noder, obj *base.Cluster) (Noder, error) {
	itemsNode, err := walk(obj.Items)
	if err != nil {
		return nil, errors.Wrap(err, "cannot handle CLUSTER.Items")
	}

	node.addAttribute("items", itemsNode)
	return node, nil
}
