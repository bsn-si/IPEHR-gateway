package base

// IsmTransition
// Model of a transition in the Instruction State Machine, caused by a careflow step.
// The attributes document the careflow step as well as the ISM transition.
// https://specifications.openehr.org/releases/RM/Release-1.0.2/ehr.html#_ism_transition_class
type IsmTransition struct {
	CurrentState DvCodedText  `json:"current_state"`
	Transition   *DvCodedText `json:"transition,omitempty"`
	CareflowStep *DvCodedText `json:"careflow_step,omitempty"`
	Reason       *[]DvText    `json:"reason,omitempty"`
	Pathable
}
