package base

// ItemTree
// Logical tree data structure. The tree may be empty.
// Used for representing data which are logically a tree such as audiology results, microbiology results, biochemistry results.
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_item_tree_class
type ItemTree struct {
	ItemStructure
	Items Items `json:"items"`
}

func (it ItemTree) GetType() ItemType {
	return ItemTreeItemType
}
