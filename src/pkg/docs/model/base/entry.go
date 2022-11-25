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
	ContentItem
	Language            CodePhrase      `json:"language"`
	Encoding            CodePhrase      `json:"encoding"`
	OtherParticipations []Participation `json:"other_participations,omitempty"`
	WorkflowID          *ObjectRef      `json:"workflow_id,omitempty"`
	Subject             PartyProxy      `json:"subject"`
	Provider            *PartyProxy     `json:"provider,omitempty"`
}

func (e Entry) GetType() ItemType {
	return e.Type
}

func (e Entry) GetLocatable() Locatable {
	return e.Locatable
}

// CareEntry
// The abstract parent of all clinical ENTRY subtypes.
// A CARE_ENTRY defines protocol and guideline attributes for all clinical Entry subtypes.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_care_entry_class
type CareEntry struct {
	Protocol    *ItemStructure `json:"protocol,omitempty"`
	GuidelineID ObjectRef      `json:"guideline_id"`
	Entry
}

// Action
// Used to record a clinical action that has been performed, which may have been ad hoc,
// or due to the execution of an Activity in an Instruction workflow.
// Every Action corresponds to a careflow step of some kind or another.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_action_class
type Action struct {
	CareEntry
	Time               DvDateTime          `json:"time"`
	IsmTransition      IsmTransition       `json:"ism_transition"`
	InstructionDetails *InstructionDetails `json:"instruction_details,omitempty"`
	Description        ItemTree            `json:"description"`
}

// Evaluation
// Entry type for evaluation statements. Used for all kinds of statements which evaluate other information,
// such as interpretations of observations, diagnoses, differential diagnoses, hypotheses,
// risk assessments,goals and plans.

// Should not be used for actionable statements such as medication orders - these are represented
// using the INSTRUCTION type.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_evaluation_class
type Evaluation struct {
	CareEntry
	Data ItemStructure `json:"data"`
}

// Instruction
// Used to specify actions in the future. Enables simple and complex specifications to be expressed, including in a fully-computable workflow form. Used for any actionable statement such as medication and therapeutic orders, monitoring, recall and review. Enough details must be provided for the specification to be directly executed by an actor, either human or machine.
//
// Not to be used for plan items which are only specified in general terms.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_observation_class
type Instruction struct {
	CareEntry
	Narrative  DvText      `json:"narrative"`
	ExpiryTime *DvDateTime `json:"expiry_time,omitempty"`
	// wf_definition *DvParsable `json:"wf_definition,omitempty`
	// activities    []Activity  `json:"activities,omitempty"`
}

// Observation
// Entry subtype for all clinical data in the past or present, i.e. which (by the time it is recorded) has already occurred.
// OBSERVATION data is expressed using the class HISTORY<T>, which guarantees that it is situated in time.
// OBSERVATION is used for all notionally objective (i.e. measured in some way) observations of phenomena,
// and patient-reported phenomena, e.g. pain.
//
// Not to be used for recording opinion or future statements of any kind, including instructions, intentions, plans etc.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_observation_class
type Observation struct {
	Data  History[ItemStructure]  `json:"data"`
	State *History[ItemStructure] `json:"state,omitempty"`
	CareEntry
}
