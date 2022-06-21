package base

// PartyIdentified
// Proxy data for an identified party other than the subject of the record, minimally consisting of
// human-readable identifier(s), such as name, formal (and possibly computable) identifiers such as NHS
// number, and an optional link to external data.
// https://specifications.openehr.org/releases/RM/latest/common.html#_party_identified_class
type PartyIdentified struct {
	Name        string         `json:"name"`
	Identifiers []DvIdentifier `json:"identifiers"`
	PartyProxy
}
