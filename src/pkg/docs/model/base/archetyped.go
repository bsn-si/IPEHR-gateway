package base

// Archetyped
// Archetypes act as the configuration basis for the particular structures of instances defined by the
// reference model.
// https://specifications.openehr.org/releases/RM/latest/common.html#_archetyped_class
type Archetyped struct {
	Type        ItemType  `json:"_type"`
	ArchetypeID ObjectID  `json:"archetype_id"`
	TemplateID  *ObjectID `json:"template_id,omitempty"` // TODO why link?
	RmVersion   string    `json:"rm_version"`
}
