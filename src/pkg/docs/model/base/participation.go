package base

// Participation
// Model of a participation of a Party (any Actor or Role) in an activity. Used to represent any
// participation of a Party in some activity, which is not explicitly in the model, e.g. assisting nurse.
// https://specifications.openehr.org/releases/RM/latest/common.html#_participation_class
type Participation struct {
	Function  DvText                `json:"function"`
	Mode      *DvCodedText          `json:"mode,omitempty"`
	Performer PartyProxy            `json:"performer"`
	Time      *Interval[DvDateTime] `json:"time,omitempty"`
}
