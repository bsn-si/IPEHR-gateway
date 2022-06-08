package base

// Evaluation
// Entry type for evaluation statements. Used for all kinds of statements which evaluate other information,
// such as interpretations of observations, diagnoses, differential diagnoses, hypotheses,
// risk assessments,goals and plans.

// Should not be used for actionable statements such as medication orders - these are represented
// using the INSTRUCTION type.
type Evaluation struct {
	Data ItemStructure `json:"data"`
	CareEntry
}
