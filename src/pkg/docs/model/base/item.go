package base

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Item
//
// The abstract parent of CLUSTER and ELEMENT representation classes.
//
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_item_class
type Item Locatable

type Items []Root

func (items *Items) UnmarshalJSON(data []byte) error {
	wrappers := []itemWrapper{}
	if err := json.Unmarshal(data, &wrappers); err != nil {
		return errors.Wrap(err, "cannot unmarshal items wrapper slice")
	}

	newItems := make(Items, 0, len(wrappers))
	for _, wrapper := range wrappers {
		newItems = append(newItems, wrapper.item)
	}

	*items = newItems

	return nil
}

type itemWrapper struct {
	item Root
}

func (itemW *itemWrapper) UnmarshalJSON(data []byte) error {
	tmpStr := struct {
		Type ItemType `json:"_type"`
	}{}
	if err := json.Unmarshal(data, &tmpStr); err != nil {
		return errors.Wrap(err, "cannot unmarshal item wrapper")
	}

	switch tmpStr.Type {
	case ElementItemType:
		itemW.item = &Element{}
	case ClusterItemType:
		itemW.item = &Cluster{}
	default:
		return errors.Errorf("unexpected item item type: '%v'", tmpStr.Type)
	}

	if err := json.Unmarshal(data, itemW.item); err != nil {
		return errors.Wrapf(err, "cannot unmarshal wrapper item type: '%v'", tmpStr.Type)
	}

	return nil
}

// Element
//
// The leaf variant of ITEM, to which a DATA_VALUE instance is attached.
//
// https://specifications.openehr.org/releases/RM/latest/data_structures.html#_element_class
type Element struct {
	Item
	NullFlavour *DvCodedText `json:"null_flavour,omitempty"`
	Value       DataValue    `json:"value,omitempty"`
	NullReason  *DvText      `json:"null_reason,omitempty"`
}

func (e Element) GetType() ItemType {
	return ElementItemType
}

func (e Element) GetLocatable() Locatable {
	return Locatable(e.Item)
}

func (e Element) GetArchetypeNodeID() string {
	return e.ArchetypeNodeID
}

func (e *Element) UnmarshalJSON(data []byte) error {
	wrapper := struct {
		Item
		NullFlavour *DvCodedText      `json:"null_flavour,omitempty"`
		Value       *dataValueWrapper `json:"value,omitempty"`
		NullReason  *DvText           `json:"null_reason,omitempty"`
	}{}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return errors.Wrap(err, "cannot unmarshal element wrapper")
	}

	e.Item = wrapper.Item
	e.NullFlavour = wrapper.NullFlavour
	e.NullFlavour = wrapper.NullFlavour

	if wrapper.Value != nil {
		e.Value = wrapper.Value.dv
	}

	return nil
}

// Cluster
//
// The grouping variant of ITEM, which may contain further instances of ITEM, in an ordered list.
//
// https://specifications.openehr.org/releases/RM/Release-1.0.2/data_structures.html#_cluster_class
type Cluster struct {
	Item
	Items Items `json:"items"`
}

func (c Cluster) GetType() ItemType {
	return ClusterItemType
}

func (c Cluster) GetLocatable() Locatable {
	return Locatable(c.Item)
}

func (c Cluster) GetArchetypeNodeID() string {
	return c.ArchetypeNodeID
}
