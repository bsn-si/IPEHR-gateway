package base

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// Abstract parent class of all spatial data types.
//
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_item_structure_class
type ItemStructure struct {
	Data Root `json:"-"`
}

func (is ItemStructure) GetType() ItemType {
	return is.Data.GetType()
}

func (is ItemStructure) GetLocatable() Locatable {
	return is.Data.GetLocatable()
}

func (is ItemStructure) GetArchetypeNodeID() string {
	return is.Data.GetArchetypeNodeID()
}

func (is ItemStructure) MarshalJSON() ([]byte, error) {
	return json.Marshal(is.Data)
}

func (is *ItemStructure) UnmarshalJSON(data []byte) error {
	tmp := struct {
		Type ItemType `json:"_type"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "cannot unmarshal item structure type")
	}

	switch tmp.Type {
	case ItemSingleItemType:
		is.Data = &ItemSingle{}
	case ItemListItemType:
		is.Data = &ItemList{}
	case ItemTableItemType:
		is.Data = &ItemTable{}
	case ItemTreeItemType:
		is.Data = &ItemTree{}
	case "":
		return nil
	default:
		return fmt.Errorf("unexpected item struct type: '%v'", tmp.Type) // nolint
	}

	if err := json.Unmarshal(data, is.Data); err != nil {
		return errors.Wrap(err, "cannot unmarshal item structure instance")
	}

	return nil
}

// ItemSingle
// Logical single value data structure.
// Used to represent any data which is logically a single value, such as a person’s height or weight.
//
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_item_single_class
type ItemSingle struct {
	DataStructure
	Item Element `json:"item"`
}

func (is ItemSingle) GetType() ItemType {
	return ItemSingleItemType
}

// ItemList
// Logical list data structure, where each item has a value and can be referred to by a name
// and a positional index in the list. The list may be empty.
//
// ITEM_LIST is used to represent any data which is logically a list of values,
// such as blood pressure, most protocols, many blood tests etc.
//
// Not to be used for time-based lists, which should be represented with the proper temporal class,
// i.e. HISTORY.
//
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_item_list_class
type ItemList struct {
	DataStructure
	Items []Element `json:"items,omitempty"`
}

func (il ItemList) GetType() ItemType {
	return ItemListItemType
}

// ItemTable
// Logical relational database style table data structure, in which columns
// are named and ordered with respect to each other. Implemented using Cluster-per-row encoding.
// Each row Cluster must have an identical number of Elements, each of which in
// turn must have identical names and value types in the corresponding positions in each row.
//
// Some columns may be designated key' columns, containing key data for each row,
// in the manner of relational tables. This allows row-naming, where each row represents
// a body site, a blood antigen etc. All values in a column have the same data type.
//
// Used for representing any data which is logically a table of values, such as blood pressure,
// most protocols, many blood tests etc.
//
// Misuse: Not to be used for time-based data, which should be represented with the
// temporal class HISTORY. The table may be empty.
//
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_item_table_class
type ItemTable struct {
	DataStructure
	Rows []Cluster `json:"rows,omitempty"`
}

func (it ItemTable) GetType() ItemType {
	return ItemTableItemType
}

// ItemTree
// Logical tree data structure. The tree may be empty.
// Used for representing data which are logically a tree such as audiology results, microbiology results, biochemistry results.
//
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_item_tree_class
type ItemTree struct {
	DataStructure
	Items Items `json:"items,omitempty"`
}

func (it ItemTree) GetType() ItemType {
	return ItemTreeItemType
}
