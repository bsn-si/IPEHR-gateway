package treeindex

import (
	"fmt"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func processItemStructure(node Noder, obj base.ItemStructure) (Noder, error) {
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

func processDataStructure(node Noder, obj *base.DataStructure) (Noder, error) {
	return node, nil
}

func processItemSingle(node Noder, obj *base.ItemSingle) (Noder, error) {
	node, err := processDataStructure(node, &obj.DataStructure)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ITEM_SINGLE.base")
	}

	itemNode, err := walk(obj.Item)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ITEM_SINGLE.item")
	}

	node.addAttribute("item", itemNode)

	return nil, errors.New("item single not implemented")
}

func processItemList(node Noder, obj *base.ItemList) (Noder, error) {
	node, err := processDataStructure(node, &obj.DataStructure)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ITEM_LIST.base")
	}

	itemsNode, err := walk(obj.Items)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ITEM_LIST.items")
	}

	node.addAttribute("items", itemsNode)

	return nil, errors.New("item list not implemented")
}

func processItemTable(node Noder, obj *base.ItemTable) (Noder, error) {
	node, err := processDataStructure(node, &obj.DataStructure)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ITEM_TABLE.base")
	}

	rowsNode, err := walk(obj.Rows)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ITEM_TABLE.rows")
	}

	node.addAttribute("rows", rowsNode)

	return nil, errors.New("item table not implemented")
}

func processItemTree(node Noder, obj *base.ItemTree) (Noder, error) {
	node, err := processDataStructure(node, &obj.DataStructure)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ITEM_TREE.base")
	}

	itemsNode, err := walk(obj.Items)
	if err != nil {
		return nil, errors.Wrap(err, "cannot process ITEM_TREE.items")
	}

	node.addAttribute("items", itemsNode)

	return node, nil
}
