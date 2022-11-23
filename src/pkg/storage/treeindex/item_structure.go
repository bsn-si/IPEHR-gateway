package treeindex

import (
	"fmt"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/errors"
)

func processItemStructure(node noder, obj base.ItemStructure) (noder, error) {
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

func processItemSingle(node noder, obj *base.ItemSingle) (noder, error) {
	return nil, errors.New("item single not implemented")
}

func processItemList(node noder, obj *base.ItemList) (noder, error) {
	fmt.Println("some text ")
	return nil, errors.New("item list not implemented")
}

func processItemTable(node noder, obj *base.ItemTable) (noder, error) {
	return nil, errors.New("item table not implemented")
}

func processItemTree(node noder, obj *base.ItemTree) (noder, error) {
	for _, item := range obj.Items {
		itemsNode, err := walk(item)
		if err != nil {
			return nil, errors.Wrap(err, "cannot process ItemTree.items")
		}

		node.addAttribute("items", itemsNode)
	}

	return node, nil
}
