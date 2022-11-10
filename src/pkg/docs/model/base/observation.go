package base

// Observation
// Entry subtype for all clinical data in the past or present, i.e. which (by the time it is recorded) has already occurred.
// OBSERVATION data is expressed using the class HISTORY<T>, which guarantees that it is situated in time.
// OBSERVATION is used for all notionally objective (i.e. measured in some way) observations of phenomena,
// and patient-reported phenomena, e.g. pain.
//
// Not to be used for recording opinion or future statements of any kind, including instructions, intentions, plans etc.
// https://specifications.openehr.org/releases/RM/latest/ehr.html#_observation_class
type Observation struct {
	Data  History[ItemStructure] `json:"data"`
	State History[ItemStructure] `json:"state,omitempty"`
	CareEntry
}
