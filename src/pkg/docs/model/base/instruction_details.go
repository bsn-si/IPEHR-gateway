package base

// InstructionDetails
// Used to record details of the Instruction causing an Action.
// https://specifications.openehr.org/releases/RM/Release-1.0.2/ehr.html#_instruction_details_class

type InstructionDetails struct {
	InstructionID LocatableRef  `json:"instruction_id"`
	ActivityID    string        `json:"activity_id"`
	WfDetails     ItemStructure `json:"wf_details,omitempty"`
	Pathable
}
