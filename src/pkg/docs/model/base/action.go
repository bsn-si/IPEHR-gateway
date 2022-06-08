package base

// Action
// Used to record a clinical action that has been performed, which may have been ad hoc,
// or due to the execution of an Activity in an Instruction workflow.
// Every Action corresponds to a careflow step of some kind or another.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_action_class
type Action struct {
	Time               DvDateTime          `json:"time"`
	IsmTransition      IsmTransition       `json:"ism_transition"`
	InstructionDetails *InstructionDetails `json:"instruction_details,omitempty"`
	Description        ItemStructure       `json:"description"`
	CareEntry
}
