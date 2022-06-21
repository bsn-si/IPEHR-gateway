package base

// Entry
// The abstract parent of all ENTRY subtypes. An ENTRY is the root of a logical item of hard clinical
// information created in the clinical statement context, within a clinical session.
// There can be numerous such contexts in a clinical session.
// Observations and other Entry types only ever document information captured/created in the
// event documented by the enclosing Composition.

// An ENTRY is also the minimal unit of information any query should return, since a whole ENTRY
// (including subparts) records spatial structure, timing information, and contextual information,
// as well as the subject and generator of the information.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_entry_class
type Entry struct {
	Language            CodePhrase       `json:"language"`
	Encoding            CodePhrase       `json:"encoding"`
	OtherParticipations *[]Participation `json:"other_participations,omitempty"`
	WorkflowId          *ObjectRef       `json:"workflow_id,omitempty"`
	Subject             PartyProxy       `json:"subject"`
	Provider            *PartyProxy      `json:"provider,omitempty"`
	ContentItem
}
