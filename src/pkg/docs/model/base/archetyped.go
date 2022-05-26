package base

// Archetyped
// Archetypes act as the configuration basis for the particular structures of instances defined by the
// reference model.
// https://specifications.openehr.org/releases/RM/latest/common.html#_archetyped_class
type Archetyped struct {
	ArchetypeId ObjectId  `json:"archetype_id"`
	TemplateId  *ObjectId `json:"template_id,omitempty"`
	RmVersion   string    `json:"rm_version"`
}
