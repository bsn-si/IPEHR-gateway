package base

// Cluster
// The grouping variant of ITEM, which may contain further instances of ITEM, in an ordered list.
// https://specifications.openehr.org/releases/RM/Release-1.0.2/data_structures.html#_cluster_class
type Cluster struct {
	Items []Item `json:"items"`
	Item
}
