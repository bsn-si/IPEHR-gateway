package base

// Party proxy representing the subject of the record.
// Used to indicate that the party is the owner of the record.
// May or may not have external_ref set.
// https://specifications.openehr.org/releases/RM/latest/common.html#_party_self_class
type PartySelf struct {
	PartyProxy
}
