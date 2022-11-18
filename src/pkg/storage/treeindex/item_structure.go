package treeindex

import (
	"errors"
	"fmt"
	"hms/gateway/pkg/docs/model/base"
)

func processItemStructure(node *Node, obj base.ItemStructure) (*Node, error) {
	switch obj.GetType() {
	case base.ItemSingleItemType:
		item, ok := obj.Data.(*base.ItemSingle)
		if !ok {
			return nil, fmt.Errorf("unexpected ItemSingle type: %T", obj.Data) // nolint
		}

		return processItemSingle(node, item)
	case base.ItemListItemType:
		item, ok := obj.Data.(*base.ItemList)
		if !ok {
			return nil, fmt.Errorf("unexpected ItemList type: %T", obj.Data) // nolint
		}

		return processItemList(node, item)
	case base.ItemTableItemType:
		item, ok := obj.Data.(*base.ItemTable)
		if !ok {
			return nil, fmt.Errorf("unexpected ItemTable type: %T", obj.Data) // nolint
		}

		return processItemTable(node, item)
	case base.ItemTreeItemType:
		item, ok := obj.Data.(*base.ItemTree)
		if !ok {
			return nil, fmt.Errorf("unexpected ItemTree type: %T", obj.Data) // nolint
		}

		return processItemTree(node, item)
	default:
		return nil, fmt.Errorf("unexpected item structure type: %v", obj.GetType()) // nolint
	}
}

func processItemSingle(node *Node, obj *base.ItemSingle) (*Node, error) {
	return nil, errors.New("not implemented")
}

func processItemList(node *Node, obj *base.ItemList) (*Node, error) {
	return nil, errors.New("not implemented")
}

func processItemTable(node *Node, obj *base.ItemTable) (*Node, error) {
	return nil, errors.New("not implemented")
}

func processItemTree(node *Node, obj *base.ItemTree) (*Node, error) {
	return nil, errors.New("not implemented")
}
