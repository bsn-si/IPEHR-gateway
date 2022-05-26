package base

// Abstract parent class of all data structure types.
// Includes the as_hierarchy function which can generate the equivalent CEN EN13606 single hierarchy
// for each subtype’s physical representation.
// For example, the physical representation of an ITEM_LIST is List<ELEMENT>;
// its implementation of as_hierarchy will generate a CLUSTER containing the set of ELEMENT nodes from the list.
//
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_data_structure_class
type DataStructure struct {
	Locatable
}
