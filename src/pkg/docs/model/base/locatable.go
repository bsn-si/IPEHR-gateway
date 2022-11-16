package base

type Root interface {
	GetType() ItemType
	GetLocatable() Locatable
	GetArchetypeNodeID() string
}

// Locatable
// Root class of all information model classes that can be archetyped. Most classes in the openEHR
// reference model inherit from the LOCATABLE class, which defines the idea of locatability in an
// archetyped structure. LOCATABLE defines a runtime name and an archetype_node_id.
// https://specifications.openehr.org/releases/RM/latest/common.html#_locatable_class
type Locatable struct {
	Type            ItemType `json:"_type"`
	Name            DvText   `json:"name"`
	ArchetypeNodeID string   `json:"archetype_node_id"`
	ObjectVersionID
	Links            *[]Link     `json:"links,omitempty"`
	ArchetypeDetails *Archetyped `json:"archetype_details,omitempty"`
	Pathable         `json:"-"`
}

func (l Locatable) GetType() ItemType {
	return l.Type
}

func (l Locatable) GetLocatable() Locatable {
	return l
}

func (l Locatable) GetArchetypeNodeID() string {
	return l.ArchetypeNodeID
}
