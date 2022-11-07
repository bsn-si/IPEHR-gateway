package base

import (
	"encoding/json"

	"github.com/pkg/errors"
)

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

type action struct {
	Time               DvDateTime          `json:"time"`
	IsmTransition      IsmTransition       `json:"ism_transition"`
	InstructionDetails *InstructionDetails `json:"instruction_details,omitempty"`
	Description        ItemStructure       `json:"description"`
	CareEntry
}

func (a *Action) UnmarshalJSON(data []byte) error {
	aa := action{}
	if err := json.Unmarshal(data, &aa); err != nil {
		return errors.Wrap(err, "cannot unmarshal action")
	}

	a.Time = aa.Time
	a.IsmTransition = aa.IsmTransition
	a.InstructionDetails = aa.InstructionDetails
	a.Description = aa.Description
	a.CareEntry = aa.CareEntry

	return nil
}
